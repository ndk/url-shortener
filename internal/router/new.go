package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"url-shortener/internal/logger"
	httplogger "url-shortener/internal/logger/http"
)

type Handlers interface {
	CreateShortLink(w http.ResponseWriter, r *http.Request)
	OpenShortLink(w http.ResponseWriter, r *http.Request)
}

func NewRouter(cfg *Config, logger *logger.Logger, handlers Handlers) http.Handler {
	r := chi.NewRouter()
	{
		r.Use(httplogger.NewHandler(*logger))
		r.Use(httplogger.Recoverer)
		r.Use(httplogger.RequestIDHandler("id_request", "X-Request-ID"))
		if !cfg.JaegerDisabled {
			r.Use(func(next http.Handler) http.Handler {
				fn := func(w http.ResponseWriter, r *http.Request) {
					span, ctx := opentracing.StartSpanFromContext(r.Context(), r.Method+" "+r.URL.String())
					defer span.Finish()
					if id, ok := httplogger.IDFromCtx(ctx); ok {
						span.LogFields(
							log.String("id_request", id.String()),
						)
					}
					r = r.WithContext(opentracing.ContextWithSpan(ctx, span))

					next.ServeHTTP(w, r)
				}
				return http.HandlerFunc(fn)
			})
		}
		if cfg.LogElapsedTime {
			r.Use(httplogger.ElapsedTime)
		}
		if cfg.LogRequests {
			r.Use(httplogger.RequestBody)
		}

		r.Post("/", handlers.CreateShortLink)
		r.Get("/{slug}", handlers.OpenShortLink)
		r.Route("/internal", func(r chi.Router) {
			r.Mount("/debug", middleware.Profiler())
		})
	}
	return r
}
