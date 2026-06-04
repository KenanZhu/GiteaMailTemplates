package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"gitea-mail-templates/tools/config"
)

// CreateCommand returns the "create" subcommand.
func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create one or more new theme style directories",
		UsageText: "go run . create [--folder <themes-dir>] [--config <config-file>] <style-name>...",
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
		},
		Action: runCreate,
	}
}

func runCreate(c *cli.Context) error {
	names := c.Args().Slice()
	if len(names) == 0 {
		return fmt.Errorf("at least one style name is required")
	}

	folder := c.String("folder")
	cfg, err := config.Load(c.String("config"))
	if err != nil {
		return err
	}

	for _, styleName := range names {
		styleDir := filepath.Join(folder, styleName)

		if info, err := os.Stat(styleDir); err == nil {
			if info.IsDir() {
				fmt.Printf("  [skip] '%s' already exists\n", styleName)
				continue
			}
			return fmt.Errorf("'%s' exists but is not a directory", styleName)
		}

		created := 0
		for _, tplCfg := range cfg.Templates {
			tmplPath := filepath.Join(styleDir, tplCfg.PathStr())
			dir := filepath.Dir(tmplPath)

			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("cannot create directory %s: %w", dir, err)
			}

			placeholder := fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head>\n"+
				"  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />\n"+
				"  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n"+
				"  <title>{{.locale.Tr \"mail.%s.title\"}}</title>\n"+
				"  <style type=\"text/css\">\n    /* TODO: Add your custom styles here */\n"+
				"  </style>\n</head>\n<body style=\"margin:0;padding:0;\">\n"+
				"  <!-- TODO: Add your custom template design for: %s -->\n"+
				"</body>\n</html>\n", tplCfg.Name, tplCfg.Desc)

			if err := os.WriteFile(tmplPath, []byte(placeholder), 0644); err != nil {
				return fmt.Errorf("cannot write %s: %w", tmplPath, err)
			}
			created++
		}

		fmt.Printf("  Created style '%s' with %d template files\n", styleName, created)
	}

	return nil
}
