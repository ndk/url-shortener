package router

type Config struct {
	JaegerDisabled bool `env:"JAEGER_DISABLED,default=false"`
	LogRequests    bool `env:"MUX_LOG_REQUESTS,default=false"`
	LogElapsedTime bool `env:"MUX_LOG_ELAPSEDTIME,default=false"`
}
