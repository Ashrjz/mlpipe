package relay

import (
	"context"
	"log/slog"

	"github.com/Ashrjz/mlpipe/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	slog.Info("relay starting")
	<-ctx.Done()
	slog.Info("relay stopping")
	return nil
}
