package service

import (
	"time"

	"url-shortener/internal/jaeger"
	"url-shortener/internal/logger"
	"url-shortener/internal/router"
	"url-shortener/internal/slugs"
	"url-shortener/internal/storage/redis"
)

type Config struct {
	Jaeger jaeger.Config
	Logger logger.Config
	Redis  redis.Config
	Router router.Config
	Slugs  slugs.Config

	Address         string        `env:"LISTEN_ADDRESS,default=:8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,default=3s"`
}
