package worker

import (
	"context"
	"log/slog"

	"github.com/Ashrjz/mlpipe/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	slog.Info("api worker")
	<-ctx.Done()
	slog.Info("api worker")
	return nil
}
