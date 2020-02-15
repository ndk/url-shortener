package jaeger

import (
	"fmt"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Setup(cfg *Config) (func() error, error) {
	jCfg := jaegercfg.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: fmt.Sprintf("%s:%s", cfg.Agent.Host, cfg.Agent.Port),
		},
	}
	closer, err := jCfg.InitGlobalTracer("")
	if err != nil {
		return nil, err
	}

	return closer.Close, nil
}
