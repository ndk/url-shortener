package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"url-shortener/internal/chi_utils"
	httplogger "url-shortener/internal/logger/http"
	"url-shortener/pkg/protocol"
)

var (
	errIncorrectSlug = errors.New("The slug is incorrect")
)

type slugsRegistry interface {
	RegisterURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, slug string) (string, error)
}

type server struct {
	slugMinLength int
	registry      slugsRegistry
	bind          func(r *http.Request, v render.Binder) error
}

func (s *server) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	request := protocol.CreateShortLinkRequest{}
	if err := s.bind(r, &request); err != nil {
		render.Render(w, r, chi_utils.InvalidRequest(err))
		return
	}

	slug, err := s.registry.RegisterURL(r.Context(), request.URL)
	if err != nil {
		httplogger.FromRequest(r).Error().Err(err).Str("url", request.URL).Msg("Cannot generate a new slug")
		render.Render(w, r, chi_utils.InternalServerError(err))
		return
	}

	response := protocol.CreateShortLinkResponse{}
	response.Data.Slug = slug
	render.Respond(w, r, &response)
}

func (s *server) OpenShortLink(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if len(slug) < s.slugMinLength {
		render.Render(w, r, chi_utils.InvalidRequest(errIncorrectSlug))
		return
	}

	url, err := s.registry.GetURL(r.Context(), slug)
	if err != nil {
		httplogger.FromRequest(r).Error().Err(err).Str("slug", slug).Msg("Cannot get an url")
		render.Render(w, r, chi_utils.InternalServerError(err))
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func NewHandlers(slugMinLength int, registry slugsRegistry) *server {
	return &server{
		slugMinLength: slugMinLength,
		registry:      registry,
		bind:          render.Bind,
	}
}
