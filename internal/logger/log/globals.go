package log

import (
	"github.com/rs/zerolog/log"

	"url-shortener/internal/logger"
)

var (
	Logger logger.Logger = log.Logger

	Fatal = log.Fatal
)
