package redis

import (
	"context"

	"github.com/go-redis/redis"
)

type storage struct {
	client           *redis.Client
	instanceIndexKey string
}

func (s *storage) Close() error {
	return s.client.Close()
}

func (s *storage) NextInstanceIndex() (int64, error) {
	return s.client.Incr(s.instanceIndexKey).Result()
}

func (s *storage) SaveValue(ctx context.Context, key string, value string) error {
	return s.client.Set(key, value, 0).Err()
}

func (s *storage) LoadValue(ctx context.Context, key string) (string, error) {
	return s.client.Get(key).Result()
}

func NewStorage(cfg *Config) *storage {
	return &storage{
		client: redis.NewClient(
			&redis.Options{
				Addr:     cfg.Address,
				Password: cfg.Password,
				DB:       cfg.Database,
			},
		),
		instanceIndexKey: cfg.InstanceIndexKey,
	}
}
