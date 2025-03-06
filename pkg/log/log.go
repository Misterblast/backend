package log

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

var Logger *slog.Logger

func InitLogger(profile string, useLoki bool, lokiURL string) {
	var lokiClient *loki.Client
	var err error

	if useLoki {
		lokiClient, err = InitLoki(lokiURL)
		if err != nil {
			fmt.Printf("Failed to initialize Loki: %v. Falling back to stdout.\n", err)
			useLoki = false
		}
	}

	SetProfileLog(profile, useLoki, lokiClient)
}

func InitLoki(uri string) (*loki.Client, error) {
	if uri == "" {
		return nil, fmt.Errorf("empty env loki_url")
	}

	config, err := loki.NewDefaultConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("error in default config: %w", err)
	}

	config.TenantID = "xyz"
	client, err := loki.New(config)
	if err != nil {
		return nil, fmt.Errorf("error initializing Loki client: %w", err)
	}

	return client, nil
}

func SetProfileLog(profile string, useLoki bool, lokiClient *loki.Client) {
	var level slog.Leveler

	switch profile {
	case "dev":
		level = slog.LevelDebug
	case "prod":
		level = slog.LevelInfo
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	if useLoki && lokiClient != nil {
		handler = slogloki.Option{Level: level, Client: lokiClient}.NewLokiHandler()
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	Logger = slog.New(handler).
		With("environment", profile).
		With("apps_name", "misterblast").
		With("release", "v0.0.1a")
}

func Debug(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Debug(msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Info(msg, args...)
	}
}

func Warn(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Warn(msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	if Logger != nil {
		Logger.Error(msg, args...)
	}
}
