package router

type Config struct {
	LogRequests    bool `env:"MUX_LOG_REQUESTS,default=false"`
	LogElapsedTime bool `env:"MUX_LOG_ELAPSEDTIME,default=false"`
}
