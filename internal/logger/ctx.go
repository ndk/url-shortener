package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func Ctx(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}
