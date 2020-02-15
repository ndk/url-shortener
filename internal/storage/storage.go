package storage

import "context"

type Storage interface {
	SaveValue(ctx context.Context, key string, value string) error
	LoadValue(ctx context.Context, key string) (string, error)
}
