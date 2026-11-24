package api

import (
	"context"
	"log/slog"

	"github.com/Ashrjz/mlpipe/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	slog.Info("api starting")
	<-ctx.Done()
	slog.Info("api stopping")
	return nil
}
