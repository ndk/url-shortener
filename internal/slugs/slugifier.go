package slugs

import (
	"errors"

	"github.com/speps/go-hashids"
)

var (
	errSlugIsCorrupted = errors.New("The slug is corrupted")
)

type hashidsSlugifier struct {
	h *hashids.HashID
}

func (s *hashidsSlugifier) NewSlug(instanceIndex int64, slugIndex int64) (string, error) {
	return s.h.EncodeInt64([]int64{instanceIndex, slugIndex})
}

func (s *hashidsSlugifier) DecodeSlug(slug string) (instanceIndex int64, slugIndex int64, err error) {
	numbers, err := s.h.DecodeInt64WithError(slug)
	if err != nil {
		return 0, 0, err
	}
	if len(numbers) != 2 {
		return 0, 0, errSlugIsCorrupted
	}
	return numbers[0], numbers[1], nil
}

func NewHashidsSlugifier(cfg *Config) (*hashidsSlugifier, error) {
	hd := hashids.NewData()
	hd.Salt = cfg.Salt
	hd.MinLength = cfg.MinLength
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	return &hashidsSlugifier{
		h: h,
	}, nil
}
