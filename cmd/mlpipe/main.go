package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ashrjz/mlpipe/internal/api"
	"github.com/Ashrjz/mlpipe/internal/config"
	"github.com/Ashrjz/mlpipe/internal/relay"
	"github.com/Ashrjz/mlpipe/internal/worker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mlpipe <api|relay|worker>")
		os.Exit(2)
	}

	var run func(context.Context, config.Config) error
	switch os.Args[1] {
	case "api":
		run = api.Run
	case "relay":
		run = relay.Run
	case "worker":
		run = worker.Run
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	if err := run(ctx, cfg); err != nil {
		slog.Error("exited with error", "err", err)
		os.Exit(1)
	}
}
