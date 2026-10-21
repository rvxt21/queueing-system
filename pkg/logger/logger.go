package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/rvxt21/queueing-system/config"
)

func NewLogger(cfg *config.Config) (*slog.Logger, error) {
	var logLvl slog.Level

	err := logLvl.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParsingLogLevel, err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLvl,
	})

	log := slog.New(handler)

	return log, nil
}
