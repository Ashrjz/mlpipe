package api

import (
	"context"
	"log/slog"

	"github.com/Ashrjz/mlpipe/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	slog.InfoContext(ctx, "api starting")
	<-ctx.Done()
	slog.InfoContext(ctx, "api stopping")
	return nil
}
