package chi_utils

import (
	"net/http"

	"github.com/go-chi/render"

	"url-shortener/pkg/protocol"
)

type errResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	protocol.ErrorResponse
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func InvalidRequest(err error) render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorResponse: protocol.ErrorResponse{
			Errors: []protocol.Error{
				protocol.Error{
					Code:        http.StatusBadRequest,
					Description: err.Error(),
				},
			},
		},
	}
}

func InternalServerError(err error) render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorResponse: protocol.ErrorResponse{
			Errors: []protocol.Error{
				protocol.Error{
					Code:        http.StatusInternalServerError,
					Description: err.Error(),
				},
			},
		},
	}
}

func NotImplementedError() render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusNotImplemented,
		ErrorResponse: protocol.ErrorResponse{
			Errors: []protocol.Error{
				protocol.Error{
					Code:        http.StatusNotImplemented,
					Description: http.StatusText(http.StatusNotImplemented),
				},
			},
		},
	}
}
