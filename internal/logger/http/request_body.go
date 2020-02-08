package httplogger

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"

	"url-shortener/internal/chi_utils"
)

//RequestBody dumps the body of the incoming request into the logger
func RequestBody(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)

		body := r.Body
		var buf []byte
		if body != nil {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				l.Error().Err(err).Msg("Cannot read out the request body")
				if err := render.Render(w, r, chi_utils.InvalidRequest(err)); err != nil {
					l.Error().Err(err).Msg("Cannot render the error")
				}
				return
			}
			buf = b
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		}

		l.Trace().Str("method", r.Method).Str("url", r.URL.String()).Interface("headers", r.Header).Bytes("request_body", buf).Msg("The incoming request")

		next.ServeHTTP(w, r)

		r.Body = body
	}
	return http.HandlerFunc(fn)
}
