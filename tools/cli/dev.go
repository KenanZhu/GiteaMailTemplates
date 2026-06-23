package cli

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"

	"gitea-mail-templates/tools/config"
	"gitea-mail-templates/tools/preview"
)

// DevCommand returns the "dev" subcommand.
func DevCommand() *cli.Command {
	return &cli.Command{
		Name:      "dev",
		Usage:     "Start a live-reload dev server with SSE browser refresh",
		UsageText: "go run . dev [--port <port>]",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 3456,
				Usage: "Port for the dev server",
			},
		},
		Action: runDev,
	}
}

func runDev(c *cli.Context) error {
	port := c.Int("port")

	// Resolve project paths relative to the tools/ directory
	toolsDir, err := os.Getwd()
	if err != nil {
		toolsDir = "."
	}
	projectRoot := filepath.Dir(toolsDir)
	themesDir := filepath.Join(projectRoot, "themes")
	previewDir := filepath.Join(projectRoot, "preview")
	configPath := filepath.Join(toolsDir, "data", "templates_config.json")

	// Verify paths exist
	for _, p := range []string{themesDir, previewDir, configPath} {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return fmt.Errorf("required path not found: %s", p)
		}
	}

	// Load template config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Start dev server
	srv := preview.NewDevServer(preview.DevServerConfig{
		Port:       port,
		ThemesDir:  themesDir,
		PreviewDir: previewDir,
		Config:     cfg,
	})

	// Handle Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n\033[32m[I]\033[0m [Server] Shutting down")
		os.Exit(0)
	}()

	return srv.Start()
}
