package logging

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

// WithCorrelationID returns a context carrying the correlation ID.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// CorrelationID returns the ID stored in ctx, or "" if there is none.
func CorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

// handler wraps another handler and adds correlation_id to every record
// whose context carries one.
type handler struct{ slog.Handler }

func (h handler) Handle(ctx context.Context, r slog.Record) error {
	if id := CorrelationID(ctx); id != "" {
		r.AddAttrs(slog.String("correlation_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs and WithGroup must re-wrap, otherwise logger.With(...) returns
// the inner handler and the correlation ID silently disappears.
func (h handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return handler{h.Handler.WithAttrs(attrs)}
}

func (h handler) WithGroup(name string) slog.Handler {
	return handler{h.Handler.WithGroup(name)}
}

// Setup installs a JSON logger on stdout as the slog default.
func Setup(service string, level slog.Level) *slog.Logger {
	base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler{base}).With("service", service)
	slog.SetDefault(logger)
	return logger
}
