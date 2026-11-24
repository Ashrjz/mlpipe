package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	APIAddr   string
	AdminAddr string
	LogLevel  slog.Level

	DatabaseURL string
	RabbitMQURL string

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool
}

// Load reads configuration from environment variables and reports every
// problem at once.
func Load() (Config, error) {
	var errs []error

	cfg := Config{
		APIAddr:   getenv("API_ADDR", ":8080"),
		AdminAddr: getenv("ADMIN_ADDR", ":8081"),
		LogLevel:  levelEnv("LOG_LEVEL", slog.LevelInfo, &errs),

		DatabaseURL: required("DATABASE_URL", &errs),
		RabbitMQURL: required("RABBITMQ_URL", &errs),

		MinIOEndpoint:  required("MINIO_ENDPOINT", &errs),
		MinIOAccessKey: required("MINIO_ACCESS_KEY", &errs),
		MinIOSecretKey: required("MINIO_SECRET_KEY", &errs),
		MinIOBucket:    getenv("MINIO_BUCKET", "datasets"),
		MinIOUseSSL:    boolEnv("MINIO_USE_SSL", false, &errs),
	}

	if err := errors.Join(errs...); err != nil {
		return Config{}, fmt.Errorf("invalid config:\n%w", err)
	}
	return cfg, nil
}

// LogValue keeps secrets out of logs when a Config is logged.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("api_addr", c.APIAddr),
		slog.String("admin_addr", c.AdminAddr),
		slog.String("log_level", c.LogLevel.String()),
		slog.String("database_url", redactURL(c.DatabaseURL)),
		slog.String("rabbitmq_url", redactURL(c.RabbitMQURL)),
		slog.String("minio_endpoint", c.MinIOEndpoint),
		slog.String("minio_bucket", c.MinIOBucket),
		slog.Bool("minio_use_ssl", c.MinIOUseSSL),
	)
}

func redactURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "***"
	}
	return u.Redacted()
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func required(key string, errs *[]error) string {
	v := os.Getenv(key)
	if v == "" {
		*errs = append(*errs, fmt.Errorf("%s is required", key))
	}
	return v
}

func boolEnv(key string, def bool, errs *[]error) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: %q is not a valid boolean", key, v))
		return def
	}
	return b
}

func levelEnv(key string, def slog.Level, errs *[]error) slog.Level {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(v)); err != nil {
		*errs = append(*errs, fmt.Errorf("%s: %q is not a valid level (debug, info, warn, error)", key, v))
		return def
	}
	return lvl
}
