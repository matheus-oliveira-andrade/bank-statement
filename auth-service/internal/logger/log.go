package logger

import (
	"log/slog"
	"os"
)

func SetupLogger() {
  jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
  })

  defaultLogger := slog.New(jsonHandler)

  slog.SetDefault(defaultLogger)
}

