package preview

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gitea-mail-templates/tools/config"
)

// DevServerConfig holds configuration for the dev server.
type DevServerConfig struct {
	Port       int
	ThemesDir  string
	PreviewDir string
	Config     *config.TemplatesConfig
}

// DevServer is a live-reload dev server with SSE-based browser notifications.
// It serves preview/ as static files, watches themes/ for .tmpl changes,
// re-renders templates in-process, and pushes reload events to connected browsers.
type DevServer struct {
	cfg      DevServerConfig
	clients  map[chan string]struct{}
	mu       sync.Mutex
	lastMods map[string]time.Time
	modsMu   sync.Mutex
}

// NewDevServer creates a new dev server instance.
func NewDevServer(cfg DevServerConfig) *DevServer {
	return &DevServer{
		cfg:      cfg,
		clients:  make(map[chan string]struct{}),
		lastMods: make(map[string]time.Time),
	}
}

// Start initializes the rendered.js, begins file watching, and serves HTTP.
// Blocks until the server stops.
func (s *DevServer) Start() error {
	// Build initial rendered.js
	if err := s.rebuildAll(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[33m[W]\033[0m [Builder] Initial render warning: %v\n", err)
	}

	// Start file watcher in background
	go s.watchThemes()

	mux := http.NewServeMux()

	// SSE endpoint — registered before static file handler
	mux.HandleFunc("/events", s.handleSSE)

	// Static file server for preview/ directory
	previewFS := http.FileServer(http.Dir(s.cfg.PreviewDir))
	mux.Handle("/", previewFS)

	addr := fmt.Sprintf(":%d", s.cfg.Port)
	fmt.Printf("\033[32m[I]\033[0m [Server] Gitea Mail Templates — Dev Server\n")
	fmt.Printf("\033[32m[I]\033[0m [Server] http://localhost%s\n", addr)
	fmt.Printf("\033[32m[I]\033[0m [Watcher] Watching themes/ for changes\n")
	fmt.Printf("\033[32m[I]\033[0m [Server] SSE endpoint ws://localhost%s/events\n\n", addr)

	return http.ListenAndServe(addr, mux)
}

// --- SSE ---

func (s *DevServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan string, 16)
	s.mu.Lock()
	s.clients[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, ch)
		s.mu.Unlock()
		close(ch)
	}()

	ctx := r.Context()
	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}

// broadcast sends a JSON message to all connected SSE clients.
func (s *DevServer) broadcast(typ string, payload map[string]interface{}) {
	msg := map[string]interface{}{"type": typ}
	for k, v := range payload {
		msg[k] = v
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.clients {
		select {
		case ch <- string(data):
		default:
			// client buffer full, skip
		}
	}
}

// --- File watcher (polling) ---

// scanMods scans themes/ for .tmpl files and returns a map of path → mod time.
func (s *DevServer) scanMods() (map[string]time.Time, error) {
	mods := make(map[string]time.Time)
	err := filepath.Walk(s.cfg.ThemesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tmpl") {
			mods[path] = info.ModTime()
		}
		return nil
	})
	return mods, err
}

// watchThemes polls the themes directory every 500ms and triggers rebuild on changes.
func (s *DevServer) watchThemes() {
	// Seed initial state
	if mods, err := s.scanMods(); err == nil {
		s.modsMu.Lock()
		s.lastMods = mods
		s.modsMu.Unlock()
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		current, err := s.scanMods()
		if err != nil {
			continue
		}

		changed := false
		s.modsMu.Lock()
		for path, mod := range current {
			last, ok := s.lastMods[path]
			if !ok || mod.After(last) {
				changed = true
			}
		}
		// Also detect deletions
		for path := range s.lastMods {
			if _, ok := current[path]; !ok {
				changed = true
			}
		}
		if changed {
			// Log which files changed
			for path, mod := range current {
				last, ok := s.lastMods[path]
				if !ok || mod.After(last) {
					rel, _ := filepath.Rel(s.cfg.ThemesDir, path)
					fmt.Printf("\033[32m[I]\033[0m [Watcher] %s changed\n", filepath.ToSlash(rel))
				}
			}
			for path := range s.lastMods {
				if _, ok := current[path]; !ok {
					rel, _ := filepath.Rel(s.cfg.ThemesDir, path)
					fmt.Printf("\033[32m[I]\033[0m [Watcher] %s deleted\n", filepath.ToSlash(rel))
				}
			}
			s.lastMods = current
			s.modsMu.Unlock()
			s.rebuildAll()
		} else {
			s.modsMu.Unlock()
		}
	}
}

// --- Rebuild ---

var rebuildMu sync.Mutex

func (s *DevServer) rebuildAll() error {
	rebuildMu.Lock()
	defer rebuildMu.Unlock()

	start := time.Now()
	fmt.Printf("\033[32m[I]\033[0m [Builder] Rebuilding all themes...\n")

	result := RenderAll(s.cfg.ThemesDir, s.cfg.Config, nil)
	if result == nil {
		err := fmt.Errorf("no results produced")
		s.broadcast("error", map[string]interface{}{"message": err.Error()})
		return err
	}

	outputPath := filepath.Join(s.cfg.PreviewDir, "rendered.js")
	if err := WriteRenderedJS(result, outputPath); err != nil {
		s.broadcast("error", map[string]interface{}{"message": err.Error()})
		return err
	}

	elapsed := time.Since(start)
	fmt.Printf("\033[32m[I]\033[0m [Builder] Rebuild done in %v\n", elapsed.Round(time.Millisecond))

	s.broadcast("reload", map[string]interface{}{"elapsed": elapsed.Milliseconds()})
	return nil
}
