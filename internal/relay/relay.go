package relay

import (
	"context"
	"log/slog"

	"github.com/Ashrjz/mlpipe/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	slog.InfoContext(ctx, "relay starting")
	<-ctx.Done()
	slog.InfoContext(ctx, "relay stopping")
	return nil
}
