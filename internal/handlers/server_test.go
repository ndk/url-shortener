package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"url-shortener/pkg/protocol"
)

type mockRegistry struct {
	m *mock.Mock
}

func (r *mockRegistry) RegisterURL(ctx context.Context, url string) (string, error) {
	args := r.m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func (r *mockRegistry) GetURL(ctx context.Context, slug string) (string, error) {
	args := r.m.Called(ctx, slug)
	return args.String(0), args.Error(1)
}

func TestCreateShortLink(t *testing.T) {
	Convey("The handler works correctly", t, func() {
		m := &mock.Mock{}
		req := httptest.NewRequest(http.MethodPost, "http://blablabla.me/", nil)
		w := httptest.NewRecorder()

		Convey("It handles the request binding errors correctly", func() {
			srv := server{
				bind: func(r *http.Request, v render.Binder) error {
					args := m.Called(r, v)
					return args.Error(0)
				},
			}
			m.
				On("1", mock.Anything, mock.Anything).Return(errors.New("Binding error"))

			srv.CreateShortLink(w, req)

			m.AssertExpectations(t)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t,
				`
        {
            "errors":
            [
                {
                    "code": 400,
                    "description": "Binding error"
                }
            ]
        }`,
				string(body),
			)
		})
		Convey("It handles the registry errors correctly", func() {
			srv := server{
				registry: &mockRegistry{
					m: m,
				},
				bind: func(r *http.Request, v render.Binder) error {
					request := v.(*protocol.CreateShortLinkRequest)
					request.URL = "http://url.me/something"
					args := m.Called(r, v)
					return args.Error(0)
				},
			}
			m.
				On("1", mock.Anything, mock.Anything).Return(nil).
				On("RegisterURL", mock.Anything, mock.Anything).Return("", errors.New("Registry error"))

			srv.CreateShortLink(w, req)

			m.AssertExpectations(t)
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t,
				`
        {
            "errors":
            [
                {
                    "code": 500,
                    "description": "Registry error"
                }
            ]
        }`,
				string(body),
			)
		})
		Convey("The handler returns a new slug", func() {
			srv := server{
				registry: &mockRegistry{
					m: m,
				},
				bind: func(r *http.Request, v render.Binder) error {
					request := v.(*protocol.CreateShortLinkRequest)
					request.URL = "http://url.me/something"
					args := m.Called(r, v)
					return args.Error(0)
				},
			}
			m.
				On("1", mock.Anything, mock.Anything).Return(nil).
				On("RegisterURL", mock.Anything, "http://url.me/something").Return("123", nil)

			srv.CreateShortLink(w, req)

			m.AssertExpectations(t)
			assert.Equal(t, http.StatusOK, w.Code)
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t,
				`
        {
            "data":
							{
									"slug": "123"
							}
        }`,
				string(body),
			)
		})
	})
}

func TestOpenShortLink(t *testing.T) {
	Convey("The handler works correctly", t, func() {
		req := httptest.NewRequest(http.MethodGet, "http://blablabla.me", nil)
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("slug", "123")
		w := httptest.NewRecorder()

		Convey("It fails if the slug is too sort", func() {
			srv := server{
				slugMinLength: 10,
			}

			srv.OpenShortLink(w, req)

			// m.AssertExpectations(t)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t,
				`
        {
            "errors":
            [
                {
                    "code": 400,
                    "description": "The slug is incorrect"
                }
            ]
        }`,
				string(body),
			)
		})
		Convey("It handles the registry errors correctly", func() {
			m := &mock.Mock{}
			srv := server{
				registry: &mockRegistry{
					m: m,
				},
				slugMinLength: 3,
			}
			m.
				On("GetURL", mock.Anything, "123").Return("", errors.New("Registry error"))

			srv.OpenShortLink(w, req)

			m.AssertExpectations(t)
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t,
				`
        {
            "errors":
            [
                {
                    "code": 500,
                    "description": "Registry error"
                }
            ]
        }`,
				string(body),
			)
		})
		Convey("It redirects to the related URL", func() {
			m := &mock.Mock{}
			srv := server{
				registry: &mockRegistry{
					m: m,
				},
				slugMinLength: 3,
			}
			m.
				On("GetURL", mock.Anything, "123").Return("http://google.com/abc", nil)

			srv.OpenShortLink(w, req)

			m.AssertExpectations(t)
			assert.Equal(t, http.StatusMovedPermanently, w.Code)
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t,
				"<a href=\"http://google.com/abc\">Moved Permanently</a>.\n\n",
				string(body),
			)
		})
	})
}

func TestNewHandlers(t *testing.T) {
	r := &mockRegistry{}
	srv := NewHandlers(73, r)
	srv.bind = nil
	assert.Equal(t,
		&server{
			slugMinLength: 73,
			registry:      r,
		},
		srv,
	)
}
