package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *Config) *Logger {
	var output io.Writer = os.Stdout
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{Out: output}
	}
	logger := zerolog.New(output).With().Logger()

	if cfg.Timestamp {
		logger = logger.With().Timestamp().Logger()
	}
	if cfg.Caller {
		logger = logger.With().Caller().Logger()
	}

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		logger.Warn().Err(err).Str("level", cfg.Level).Msg("Cannot parse the logging level")
	} else {
		logger = logger.Level(level)
	}

	return &logger
}
