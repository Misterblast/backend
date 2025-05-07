package middleware

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

var Logger *slog.Logger

type customTextHandler struct {
	writer io.Writer
	level  slog.Leveler
	attrs  []slog.Attr
	group  string
}

func extractRecordAttrs(r slog.Record) []slog.Attr {
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	return attrs
}

func (h *customTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := append([]slog.Attr{}, h.attrs...)
	merged = append(merged, attrs...)
	return &customTextHandler{
		writer: h.writer,
		level:  h.level,
		attrs:  merged,
		group:  h.group,
	}
}

func (h *customTextHandler) WithGroup(name string) slog.Handler {
	return &customTextHandler{
		writer: h.writer,
		level:  h.level,
		attrs:  h.attrs,
		group:  name,
	}
}

func NewCustomTextHandler(w io.Writer, level slog.Leveler) slog.Handler {
	return &customTextHandler{writer: w, level: level}
}

func (h *customTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *customTextHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := r.Time.Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("%s [%s] %s", timestamp, r.Level.String(), r.Message)

	// Render attributes dari handler (global) dan record (local)
	allAttrs := append(h.attrs, extractRecordAttrs(r)...)
	for _, attr := range allAttrs {
		logLine += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
	}

	if h.group != "" {
		logLine = fmt.Sprintf("[%s] %s", h.group, logLine)
	}

	logLine += "\n"
	_, err := h.writer.Write([]byte(logLine))
	return err
}

func Log(profile string, useLoki bool, lokiURL string) {
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
	case "stg":
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
		handler = NewCustomTextHandler(os.Stdout, level)
	}

	Logger = slog.New(handler)
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
