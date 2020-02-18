package httplogger

import (
	"github.com/rs/zerolog/hlog"
)

var NewHandler = hlog.NewHandler
var IDFromCtx = hlog.IDFromCtx
var RequestIDHandler = hlog.RequestIDHandler
var FromRequest = hlog.FromRequest
