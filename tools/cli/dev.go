package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"
)

// DevCommand returns the "dev" subcommand.
func DevCommand() *cli.Command {
	return &cli.Command{
		Name:  "dev",
		Usage: "Start a live-reload dev server with CSS inlining (requires Node.js)",
		UsageText: `go run . dev [--port <port>]`,
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

	// Find the server directory relative to the tools/ dir
	exe, _ := os.Executable()
	toolsDir := filepath.Dir(exe)
	// When run with `go run`, resolve relative to CWD
	if _, err := os.Stat(filepath.Join(toolsDir, "server")); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		toolsDir = cwd
	}
	serverDir := filepath.Join(toolsDir, "server")

	// Verify Node.js is available
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("Node.js is required but not found in PATH")
	}

	// Verify node_modules are installed
	if _, err := os.Stat(filepath.Join(serverDir, "node_modules")); os.IsNotExist(err) {
		fmt.Println("Installing Node.js dependencies...")
		cmd := exec.Command("npm", "install")
		cmd.Dir = serverDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("npm install failed: %w", err)
		}
	}

	// Ensure rendered.js exists before starting
	initCmd := exec.Command("go", "run", ".", "preview", "all")
	initCmd.Dir = toolsDir
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	if err := initCmd.Run(); err != nil {
		fmt.Println("Warning: initial preview generation failed, starting anyway...")
	}

	// Start the Node.js dev server
	serverCmd := exec.Command("node", "server.mjs")
	serverCmd.Dir = serverDir
	serverCmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start dev server: %w", err)
	}

	fmt.Printf("\nDev server running at http://localhost:%d\n", port)
	fmt.Print("Press Ctrl+C to stop.\n\n")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		if err := serverCmd.Process.Signal(os.Interrupt); err != nil {
			serverCmd.Process.Kill()
		}
		os.Exit(0)
	}()

	if err := serverCmd.Wait(); err != nil {
		return fmt.Errorf("dev server exited: %w", err)
	}

	return nil
}
