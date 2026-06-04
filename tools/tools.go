package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	cmds "gitea-mail-templates/tools/cli"
)

func main() {
	app := &cli.App{
		Name:      "tools",
		Usage:     "Gitea Mail Templates CLI — manage and preview email template themes",
		UsageText: "tools <command> [args...] [--flags]",
		Commands: []*cli.Command{
			cmds.ListCommand(),
			cmds.CreateCommand(),
			cmds.DeleteCommand(),
			cmds.PreviewCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
