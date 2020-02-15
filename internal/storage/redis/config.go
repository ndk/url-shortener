package redis

type Config struct {
	Address          string `env:"REDIS_ADDRESS,required"`
	Database         int    `env:"REDIS_DATABASE,required"`
	Password         string `env:"REDIS_PASSWORD"`
	InstanceIndexKey string `env:"REDIS_INSTANCEINDEXKEY,default=instance_index"`
}
