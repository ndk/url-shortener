package httplogger

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/render"

	"url-shortener/internal/chi_utils"
)

//ErrPanic We don't want to show the internal information about how our service works, so let's show it like the internal error
var ErrPanic = errors.New(http.StatusText(http.StatusInternalServerError))

//Recoverer suppresses panics
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if info := recover(); info != nil {
				FromRequest(r).Error().Interface("recover_info", info).Bytes("debug_stack", debug.Stack()).Msg("panic_on_request")
				if err := render.Render(w, r, chi_utils.InternalServerError(ErrPanic)); err != nil {
					FromRequest(r).Error().Msg("Cannot render the error")
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
