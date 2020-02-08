package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis"
	"github.com/joeshaw/envdecode"
	"github.com/oklog/run"

	"url-shortener/internal/handlers"
	"url-shortener/internal/logger"
	"url-shortener/internal/logger/log"
	"url-shortener/internal/router"
	"url-shortener/internal/slugs"
)

type config struct {
	Logger logger.Config
	Router router.Config
	Slugs  slugs.Config

	Address string `env:"LISTEN_ADDRESS,default=:8080"`
	Redis   struct {
		Address          string `env:"REDIS_ADDRESS,required"`
		Database         int    `env:"REDIS_DATABASE,required"`
		Password         string `env:"REDIS_PASSWORD"`
		InstanceIndexKey string `env:"REDIS_INSTANCEINDEXKEY,default=instance_index"`
	}
}

var (
	gitCommit = "undefined"
	gitBranch = "undefined"
)

func main() {
	cfg := &config{}
	if err := envdecode.StrictDecode(cfg); err != nil {
		log.Fatal().Err(err).Str("git_commit", gitCommit).Str("git_branch", gitBranch).Msg("Cannot decode the envs to the config")
	}

	l := logger.NewLogger(&cfg.Logger)
	l.Info().Str("git_commit", gitCommit).Str("git_branch", gitBranch).Interface("config", cfg).Msg("The gathered config")

	ctx := l.WithContext(context.Background())
	ctx, cancel := context.WithCancel(l.WithContext(ctx))

	/////////////////////////////////////////////////////////////////////////////

	storage := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Database,
		},
	)
	defer func() {
		if err := storage.Close(); err != nil {
			l.Error().Err(err).Msg("The storage has been closed improperly")
		}
	}()

	instanceIndex, err := storage.Incr(cfg.Redis.InstanceIndexKey).Result()
	if err != nil {
		l.Fatal().Err(err).Msg("Cannot retrieve the service instance index")
	}

	/////////////////////////////////////////////////////////////////////////////

	g := &run.Group{}
	{
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			cancel()
			close(stop)
		})
	}
	{
		slugifier, err := slugs.NewHashidsSlugifier(&cfg.Slugs)
		if err != nil {
			l.Fatal().Err(err).Msg("Cannot create a new slugifier")
		}
		registry := slugs.NewRegistry(
			slugifier,
			func(key string, value string) error {
				return storage.Set(key, value, 0).Err()
			},
			func(key string) (string, error) {
				return storage.Get(key).Result()
			},
			instanceIndex,
		)
		h := handlers.NewHandlers(cfg.Slugs.MinLength, registry)
		r := router.NewRouter(&cfg.Router, l, h)
		srv := http.Server{
			Addr:    cfg.Address,
			Handler: r,
		}

		g.Add(func() error {
			l.Info().Str("address", srv.Addr).Msg("Start listening")
			if err := srv.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					return nil
				}
				return err
			}
			l.Info().Msg("Listening has been stopped")
			return nil
		}, func(err error) {
			l.Info().Err(err).Msg("Shutting down of listening...")
			if err := srv.Shutdown(ctx); err != nil {
				l.Error().Err(err).Msg("Cannot shut down listening properly")
			}
			l.Info().Msg("The service has been shut down")
		})
	}

	/////////////////////////////////////////////////////////////////////////////

	l.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		l.Fatal().Err(err).Msg("The service has been stopped with the error")
	}
	l.Info().Msg("The service has been stopped")
}
