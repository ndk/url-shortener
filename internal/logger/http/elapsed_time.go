package httplogger

import (
	"net/http"
	"time"
)

//ElapsedTime reports into the logger how much time was taken to process the request
func ElapsedTime(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		stop := time.Now()
		FromRequest(r).Trace().
			Str("start_time", start.Format(time.RFC3339Nano)).
			Str("stop_time", stop.Format(time.RFC3339Nano)).
			TimeDiff("duration", stop, start).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Msg("The elapsed time")
	}
	return http.HandlerFunc(fn)
}
