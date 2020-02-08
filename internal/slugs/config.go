package slugs

type Config struct {
	Salt      string `env:"SLUGS_SALT,required"`
	MinLength int    `env:"SLUGS_MINLENGTH,default=30"`
}
