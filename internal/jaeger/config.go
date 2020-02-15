package jaeger

type Config struct {
	Agent struct {
		Host string `env:"JAEGER_AGENT_HOST,required"`
		Port string `env:"JAEGER_AGENT_PORT,required"`
	}
	Disabled    bool   `env:"JAEGER_DISABLED,default=false"`
	ServiceName string `env:"JAEGER_SERVICE_NAME,default=url-shortener"`
}
