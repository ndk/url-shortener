package storage

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type otStorage struct {
	storage Storage
}

func (s *otStorage) SaveValue(ctx context.Context, key string, value string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SaveValue")
	defer span.Finish()
	return s.storage.SaveValue(ctx, key, value)
}

func (s *otStorage) LoadValue(ctx context.Context, key string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadValue")
	defer span.Finish()
	return s.storage.LoadValue(ctx, key)
}

func TraceStorage(storage Storage) Storage {
	return &otStorage{
		storage: storage,
	}
}
