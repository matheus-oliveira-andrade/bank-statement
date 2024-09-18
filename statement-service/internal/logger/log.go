package logger

import (
	"log/slog"
	"os"
)

func SetupLogger(serviceName string) {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	defaultLogger := slog.New(jsonHandler)
	defaultLogger = defaultLogger.With("service", serviceName)

	slog.SetDefault(defaultLogger)
}
