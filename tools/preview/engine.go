package preview

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitea-mail-templates/tools/config"
)

// RenderResult holds the rendered HTML for a single template in a single theme.
type RenderResult struct {
	HTML    string
	Size    int64 // bytes of the rendered HTML
	Error   string
	Skipped bool // true when the template file was not found
}

// ThemeSummary collects statistics for one theme directory.
type ThemeSummary struct {
	Name           string
	TotalFiles     int
	TotalSize      int64 // total size of .tmpl source files
	RenderedSize   int64 // total size of rendered HTML
	RenderedCount  int
	MissingCount   int
	ErrorCount     int
	RenderDuration time.Duration
	Missing        []string                 // template IDs that were not found
	Errors         []string                 // template IDs that had render errors
	Results        map[string]*RenderResult // per-template results
}

// PreviewResult is the full output from a preview run.
type PreviewResult struct {
	Rendered  map[string]map[string]string `json:"rendered"`
	Registry  map[string]map[string]string `json:"registry"`
	Params    map[string]map[string]string `json:"params"`
	Summaries []ThemeSummary               `json:"-"`
}

// DiscoverThemes scans a directory and returns all subdirectory names.
func DiscoverThemes(themesDir string) ([]string, error) {
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read themes directory %s: %w", themesDir, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no theme directories found in %s", themesDir)
	}
	return names, nil
}

// RenderAll renders all templates from the config for all discovered themes.
// If themeFilter is non-empty, only themes in the filter set are rendered.
func RenderAll(themesDir string, cfg *config.TemplatesConfig, themeFilter map[string]bool) *PreviewResult {
	themes, err := DiscoverThemes(themesDir)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}

	// Apply theme filter
	if len(themeFilter) > 0 {
		filtered := themes[:0]
		for _, t := range themes {
			if themeFilter[t] {
				filtered = append(filtered, t)
			}
		}
		themes = filtered
		if len(themes) == 0 {
			log.Printf("WARN: no themes matched the requested filter")
			return nil
		}
	}

	result := &PreviewResult{
		Rendered: make(map[string]map[string]string),
		Registry: config.Registry(cfg),
	}
	// Build flattened params for the preview panel
	result.Params = make(map[string]map[string]string, len(cfg.Templates))
	for id, t := range cfg.Templates {
		result.Params[id] = config.FlattenParams(t.Params)
	}

	for _, themeName := range themes {
		themeDir := filepath.Join(themesDir, themeName)
		summary := renderTheme(themeDir, themeName, cfg)
		result.Summaries = append(result.Summaries, summary)

		if result.Rendered[themeName] == nil {
			result.Rendered[themeName] = make(map[string]string)
		}
		for tplID, rr := range summary.Results {
			if rr.Error != "" {
				result.Rendered[themeName][tplID] = fmt.Sprintf(
					"<p style='color:red'>Render error for %s: %s</p>", tplID, rr.Error)
			} else if rr.Skipped {
				result.Rendered[themeName][tplID] = fmt.Sprintf(
					"<p style='color:red'>Template not found: %s</p>", tplID)
			} else {
				result.Rendered[themeName][tplID] = rr.HTML
			}
		}
	}

	return result
}

func renderTheme(themeDir, themeName string, cfg *config.TemplatesConfig) ThemeSummary {
	summary := ThemeSummary{
		Name:    themeName,
		Results: make(map[string]*RenderResult),
	}

	start := time.Now()

	for tplID, tplCfg := range cfg.Templates {
		tmplRelPath := tplCfg.PathStr()
		tmplAbsPath := filepath.Join(themeDir, tmplRelPath)

		// Check source file existence
		info, err := os.Stat(tmplAbsPath)
		if os.IsNotExist(err) {
			summary.MissingCount++
			summary.Missing = append(summary.Missing, tplID)
			summary.Results[tplID] = &RenderResult{Skipped: true}
			continue
		}
		if err != nil {
			summary.ErrorCount++
			summary.Errors = append(summary.Errors, tplID)
			summary.Results[tplID] = &RenderResult{Error: err.Error()}
			continue
		}
		summary.TotalFiles++
		summary.TotalSize += info.Size()

		// Read template content
		tmplContent, err := os.ReadFile(tmplAbsPath)
		if err != nil {
			summary.ErrorCount++
			summary.Errors = append(summary.Errors, tplID)
			summary.Results[tplID] = &RenderResult{Error: fmt.Sprintf("read error: %v", err)}
			continue
		}

		// Build data context
		data := make(map[string]any)
		for k, v := range tplCfg.Params {
			data[k] = v
		}
		data["locale"] = Locale{}

		// Parse and execute
		tmpl, err := template.New(tplID).Funcs(TemplateFuncs).Parse(string(tmplContent))
		if err != nil {
			summary.ErrorCount++
			summary.Errors = append(summary.Errors, tplID)
			summary.Results[tplID] = &RenderResult{Error: fmt.Sprintf("parse error: %v", err)}
			continue
		}

		var buf strings.Builder
		if err := tmpl.Execute(&buf, data); err != nil {
			summary.ErrorCount++
			summary.Errors = append(summary.Errors, tplID)
			summary.Results[tplID] = &RenderResult{Error: fmt.Sprintf("render error: %v", err)}
			continue
		}

		rendered := buf.String()

		summary.RenderedCount++
		summary.RenderedSize += int64(len(rendered))
		summary.Results[tplID] = &RenderResult{
			HTML: rendered,
			Size: int64(len(rendered)),
		}
	}

	summary.RenderDuration = time.Since(start)
	return summary
}

// WriteRenderedJS writes the preview result to a JavaScript file loaded by preview/index.html.
func WriteRenderedJS(result *PreviewResult, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", dir, err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file %s: %w", outputPath, err)
	}
	defer outFile.Close()

	outFile.WriteString("// Auto-generated by tools/tools.go — do not edit\n")

	// Write __RENDERED__
	renderedJSON, err := json.MarshalIndent(result.Rendered, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal rendered data: %w", err)
	}
	outFile.WriteString("window.__RENDERED__ = ")
	outFile.Write(renderedJSON)
	outFile.WriteString(";\n\n")

	// Write __REGISTRY__
	registryJSON, err := json.MarshalIndent(result.Registry, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal registry data: %w", err)
	}
	outFile.WriteString("window.__REGISTRY__ = ")
	outFile.Write(registryJSON)
	outFile.WriteString(";\n\n")

	// Write __PARAMS__
	paramsJSON, err := json.MarshalIndent(result.Params, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal params data: %w", err)
	}
	outFile.WriteString("window.__PARAMS__ = ")
	outFile.Write(paramsJSON)
	outFile.WriteString(";\n")

	return nil
}

// PrintDetailedSummary prints a human-readable summary of the preview results.
func PrintDetailedSummary(result *PreviewResult, themesDir string, cfg *config.TemplatesConfig) {
	absDir, _ := filepath.Abs(themesDir)
	fmt.Printf("\nfound %d styles of possible Gitea mail template in '%s' folder:\n\n",
		len(result.Summaries), filepath.Base(absDir))

	for _, s := range result.Summaries {
		if s.MissingCount == 0 && s.ErrorCount == 0 {
			fmt.Printf("    %-16s all required mail templates are properly rendered, (total %d .tmpl files, %s after, %v)\n",
				s.Name, s.TotalFiles,
				formatSize(s.RenderedSize),
				s.RenderDuration.Round(time.Millisecond))
		} else {
			fmt.Printf("    %-16s only %d required mail templates are properly rendered:\n",
				s.Name, s.RenderedCount)
			for _, tplID := range s.Missing {
				if t, ok := cfg.Templates[tplID]; ok {
					fmt.Printf("        [%s] not found, expected at 'themes/%s/%s'\n",
						filepath.Base(t.PathStr()), s.Name, t.PathStr())
				} else {
					fmt.Printf("        [%s] not found\n", tplID)
				}
			}
			for _, tplID := range s.Errors {
				if rr, ok := s.Results[tplID]; ok && rr.Error != "" {
					fmt.Printf("        [%s] render error: %s\n", tplID, rr.Error)
				} else {
					fmt.Printf("        [%s] render error\n", tplID)
				}
			}
		}
	}
	fmt.Println()
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	return fmt.Sprintf("%dKiB", bytes/unit)
}
