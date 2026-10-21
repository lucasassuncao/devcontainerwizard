package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lucasassuncao/bezel/theme"
	"github.com/lucasassuncao/devcontainerwizard/internal/cmd"
	"github.com/lucasassuncao/devcontainerwizard/internal/updater"
)

// version is set at build time via -ldflags "-X main.version=<tag>"
var version = "dev"

func main() {
	if err := run(); err != nil {
		th := theme.Resolve(theme.Terminal, true)
		fmt.Fprintln(os.Stderr, th.Danger.Render("Error:"), err)
		os.Exit(1)
	}
}

// run wires up the application and executes it. It exists so the signal
// handler's cleanup runs on every path: main calls os.Exit on failure, which
// would skip any defer declared alongside it.
func run() error {
	updater.CleanOldBinary()

	// Ctrl+C (and SIGTERM) cancels the context instead of killing the process
	// outright, so a command in progress can stop cleanly rather than leaving a
	// half-written devcontainer.json behind.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return cmd.RootCmd(version).ExecuteContext(ctx)
}
