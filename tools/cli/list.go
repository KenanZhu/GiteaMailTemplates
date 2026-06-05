package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/urfave/cli/v2"

	"gitea-mail-templates/tools/preview"
)

// ListCommand returns the "list" subcommand.
func ListCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List available theme styles in the target folder",
		UsageText: "go run . list [--folder <themes-dir>]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "folder",
				Value: "../themes",
				Usage: "Path to the themes directory",
			},
		},
		Action: runList,
	}
}

func runList(c *cli.Context) error {
	folder := c.String("folder")

	themes, err := preview.DiscoverThemes(folder)
	if err != nil {
		return err
	}

	absDir, _ := filepath.Abs(folder)
	fmt.Printf("\033[32m[I]\033[0m [Preview] Available styles in '%s' (%d):\n", filepath.Base(absDir), len(themes))

	sort.Strings(themes)
	for _, t := range themes {
		count := 0
		filepath.Walk(filepath.Join(folder, t), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if filepath.Ext(path) == ".tmpl" {
				count++
			}
			return nil
		})
		fmt.Printf("\033[32m[I]\033[0m [Preview]    %-16s (%d .tmpl files)\n", t, count)
	}
	return nil
}
