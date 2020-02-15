package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"

	"url-shortener/internal/handlers"
	"url-shortener/internal/jaeger"
	"url-shortener/internal/logger"
	"url-shortener/internal/router"
	"url-shortener/internal/slugs"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/redis"
)

func Run(cfg *Config, l *logger.Logger) error {
	g := &run.Group{}

	{
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			close(stop)
		})
	}

	{
		redis := redis.NewStorage(&cfg.Redis)
		defer func() {
			if err := redis.Close(); err != nil {
				l.Error().Err(err).Msg("The storage has been closed improperly")
			}
		}()

		instanceIndex, err := redis.NextInstanceIndex()
		if err != nil {
			l.Error().Err(err).Msg("Cannot retrieve the service instance index")
			return err
		}

		var s storage.Storage = redis

		/////////////////////////////////////////////////////////////////////////////
		if !cfg.Jaeger.Disabled {
			close, err := jaeger.Setup(&cfg.Jaeger)
			if err != nil {
				l.Error().Err(err).Msg("Couldn't setup jaeger tracer")
				return err
			}
			defer func() {
				if err := close(); err != nil {
					l.Error().Err(err).Msg("Couldn't close jaeger tracer")
				}
			}()

			s = storage.TraceStorage(s)
		}

		slugifier, err := slugs.NewHashidsSlugifier(&cfg.Slugs)
		if err != nil {
			l.Error().Err(err).Msg("Cannot create a new slugifier")
			return err
		}
		registry := slugs.NewRegistry(slugifier, s, instanceIndex)
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

			ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
			ctx = l.WithContext(ctx)
			defer cancel()

			srv.SetKeepAlivesEnabled(false)
			if err := srv.Shutdown(ctx); err != nil {
				l.Error().Err(err).Msg("Cannot shut down listening properly")
			}
			l.Info().Msg("The service has been shut down")
		})
	}

	return g.Run()
}
