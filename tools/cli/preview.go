package cli

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"gitea-mail-templates/tools/config"
	"gitea-mail-templates/tools/preview"
)

// PreviewCommand returns the "preview" subcommand.
func PreviewCommand() *cli.Command {
	return &cli.Command{
		Name:  "preview",
		Usage: "Pre-render mail templates for preview/testing",
		UsageText: `go run . preview [--folder <themes-dir>] [--config <config-file>] [--output <output.js>] all
   go run . preview [--folder <themes-dir>] [--config <config-file>] <style-name>...`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "folder",
				Value: "../themes",
				Usage: "Path to the themes directory",
			},
			&cli.StringFlag{
				Name:  "config",
				Value: "./data/templates_config.json",
				Usage: "Path to templates_config.json",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Path to output rendered.js (default: derived from --folder)",
			},
		},
		Action: runPreview,
	}
}

func runPreview(c *cli.Context) error {
	folder := c.String("folder")
	configPath := c.String("config")

	outputPath := c.String("output")
	if outputPath == "" {
		outputPath = filepath.Join(filepath.Dir(folder), "preview", "rendered.js")
	}

	names := c.Args().Slice()
	if len(names) == 0 {
		return fmt.Errorf("specify style names or 'all' (e.g. go run . preview all)")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	var themeFilter map[string]bool
	if len(names) == 1 && names[0] == "all" {
		themeFilter = nil
	} else {
		themeFilter = make(map[string]bool)
		for _, name := range names {
			themeFilter[name] = true
		}
	}

	result := preview.RenderAll(folder, cfg, themeFilter)
	if result == nil {
		return fmt.Errorf("no results produced")
	}

	if err := preview.WriteRenderedJS(result, outputPath); err != nil {
		return err
	}

	preview.PrintDetailedSummary(result, folder, cfg)

	log.Printf("Wrote %d themes x %d templates to %s",
		len(result.Summaries), len(cfg.Templates), outputPath)

	return nil
}
