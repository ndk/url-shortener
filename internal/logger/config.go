package logger

type Config struct {
	Level     string `env:"LOGGER_LEVEL,default=info"`
	Timestamp bool   `env:"LOGGER_TIMESTAMP,default=false"`
	Caller    bool   `env:"LOGGER_CALLER,default=true"`
	Pretty    bool   `env:"LOGGER_PRETTY,default=false"`
}
