package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// DeleteCommand returns the "delete" subcommand.
func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete one or more theme style directories",
		UsageText: "go run . delete [--folder <themes-dir>] <style-name>...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "folder",
				Value: "../themes",
				Usage: "Path to the themes directory",
			},
		},
		Action: runDelete,
	}
}

func runDelete(c *cli.Context) error {
	names := c.Args().Slice()
	if len(names) == 0 {
		return fmt.Errorf("at least one style name is required")
	}

	folder := c.String("folder")

	for _, styleName := range names {
		styleDir := filepath.Join(folder, styleName)

		if _, err := os.Stat(styleDir); os.IsNotExist(err) {
			fmt.Printf("\033[33m[W]\033[0m [CLI] '%s' does not exist, skipped\n", styleName)
			continue
		}

		if err := os.RemoveAll(styleDir); err != nil {
			return fmt.Errorf("cannot delete %s: %w", styleDir, err)
		}
		fmt.Printf("\033[32m[I]\033[0m [CLI] Deleted style '%s'\n", styleName)
	}

	return nil
}
