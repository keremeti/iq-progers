package logger

import (
	"os"

	"log/slog"

	"github.com/keremeti/iq-progers/config"
)

type Logger struct {
	*slog.Logger
}

func New(env config.Configuration) *Logger {
	var log *slog.Logger

	switch env {
	case config.Dev:
		log = setupPrettySlog()
	case config.Test:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.Release:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return &Logger{log}
}

func setupPrettySlog() *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
