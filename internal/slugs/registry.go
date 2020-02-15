package slugs

import (
	"context"
	"fmt"

	"url-shortener/internal/logger"
	"url-shortener/internal/storage"
)

type slugifier interface {
	NewSlug(instanceIndex int64, slugIndex int64) (string, error)
	DecodeSlug(slug string) (instanceIndex int64, slugIndex int64, err error)
}

type registry struct {
	slugifier     slugifier
	storage       storage.Storage
	instanceIndex int64
	slugsCount    int64
}

func (r *registry) RegisterURL(ctx context.Context, url string) (string, error) {
	slug, err := r.slugifier.NewSlug(r.instanceIndex, r.slugsCount)
	if err != nil {
		return "", err
	}
	logger.Ctx(ctx).Trace().Str("slug", slug).Msg("The new slug has been produced")

	key := fmt.Sprintf("%d:%d", r.instanceIndex, r.slugsCount)
	if err := r.storage.SaveValue(ctx, key, url); err != nil {
		logger.Ctx(ctx).Error().Err(err).Str("key", key).Str("url", url).Msg("Cannot create a record")
		return "", err
	}

	r.slugsCount++
	return slug, nil
}

func (r *registry) GetURL(ctx context.Context, slug string) (string, error) {
	instanceIndex, slugIndex, err := r.slugifier.DecodeSlug(slug)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%d:%d", instanceIndex, slugIndex)
	url, err := r.storage.LoadValue(ctx, key)
	if err != nil {
		logger.Ctx(ctx).Error().Err(err).Str("key", key).Msg("Cannot read a value")
		return "", err
	}

	return url, nil
}

func NewRegistry(slugifier slugifier, storage storage.Storage, instanceIndex int64) *registry {
	return &registry{
		slugifier:     slugifier,
		storage:       storage,
		instanceIndex: instanceIndex,
	}
}
