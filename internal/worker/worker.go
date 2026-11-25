package worker

import (
	"context"
	"log/slog"

	"github.com/Ashrjz/mlpipe/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	slog.InfoContext(ctx, "api worker")
	<-ctx.Done()
	slog.InfoContext(ctx, "api worker")
	return nil
}
