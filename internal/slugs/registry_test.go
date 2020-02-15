package slugs

import (
	"context"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStorage struct {
	m *mock.Mock
}

func (s *mockStorage) SaveValue(ctx context.Context, key string, value string) error {
	args := s.m.Called(ctx, key, value)
	return args.Error(0)
}

func (s *mockStorage) LoadValue(ctx context.Context, key string) (string, error) {
	args := s.m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

type mockSlugifier struct {
	m *mock.Mock
}

func (s *mockSlugifier) NewSlug(instanceIndex int64, slugIndex int64) (string, error) {
	args := s.m.Called(instanceIndex, slugIndex)
	return args.String(0), args.Error(1)
}

func (s *mockSlugifier) DecodeSlug(slug string) (instanceIndex int64, slugIndex int64, err error) {
	args := s.m.Called(slug)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

func TestRegisterURL(t *testing.T) {
	Convey("Test RegisterURL", t, func() {
		m := &mock.Mock{}

		r := registry{
			slugifier:     &mockSlugifier{m: m},
			instanceIndex: 5,
			slugsCount:    19,
			storage:       &mockStorage{m: m},
		}

		Convey("It fails if the slugifier has failed", func() {
			m.
				On("NewSlug", int64(5), int64(19)).Return("", errors.New("NewSlug error"))

			_, err := r.RegisterURL(nil, "")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "NewSlug error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It fails if the value cannot be saved", func() {
			m.
				On("NewSlug", int64(5), int64(19)).Return("qwe", nil).
				On("SaveValue", mock.Anything, "5:19", "http://en.wikipedia.com").Return(errors.New("saveValue error"))

			_, err := r.RegisterURL(context.TODO(), "http://en.wikipedia.com")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "saveValue error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It returns a new slug", func() {
			m.
				On("NewSlug", int64(5), int64(19)).Return("qwe", nil).
				On("SaveValue", mock.Anything, "5:19", "http://en.wikipedia.com").Return(nil).
				On("NewSlug", int64(5), int64(20)).Return("asd", nil).
				On("SaveValue", mock.Anything, "5:20", "http://en.wikipedia.com").Return(nil)

			{
				slug, err := r.RegisterURL(context.TODO(), "http://en.wikipedia.com")
				assert.NoError(t, err)
				assert.Equal(t, "qwe", slug)
				assert.Equal(t, int64(5), r.instanceIndex)
				assert.Equal(t, int64(20), r.slugsCount)
			}
			{
				slug, err := r.RegisterURL(context.TODO(), "http://en.wikipedia.com")
				assert.NoError(t, err)
				assert.Equal(t, "asd", slug)
				assert.Equal(t, int64(5), r.instanceIndex)
				assert.Equal(t, int64(21), r.slugsCount)
			}

			m.AssertExpectations(t)
		})
	})
}

func TestGetURL(t *testing.T) {
	Convey("Test GetURL", t, func() {
		m := &mock.Mock{}

		r := registry{
			slugifier:     &mockSlugifier{m: m},
			instanceIndex: 5,
			slugsCount:    19,
			storage:       &mockStorage{m: m},
		}

		Convey("It fails if the slugifier has failed", func() {
			m.
				On("DecodeSlug", "123").Return(int64(321), int64(432), errors.New("DecodeSlug error"))

			_, err := r.GetURL(context.TODO(), "123")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "DecodeSlug error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It fails if the value cannot be loaded", func() {
			m.
				On("DecodeSlug", "123").Return(int64(321), int64(432), nil).
				On("LoadValue", mock.Anything, "321:432").Return("", errors.New("loadValue error"))

			_, err := r.GetURL(context.TODO(), "123")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "loadValue error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It returns the correct URL", func() {
			m.
				On("DecodeSlug", "123").Return(int64(321), int64(432), nil).
				On("LoadValue", mock.Anything, "321:432").Return("http://uber.com", nil)

			url, err := r.GetURL(context.TODO(), "123")

			m.AssertExpectations(t)
			assert.NoError(t, err)
			assert.Equal(t, "http://uber.com", url)
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})
	})
}

func TestNewRegistry(t *testing.T) {
	slugifier := &mockSlugifier{}
	storage := &mockStorage{}
	r := NewRegistry(slugifier, storage, 178)
	assert.Equal(t,
		&registry{
			slugifier:     slugifier,
			storage:       storage,
			instanceIndex: 178,
		},
		r,
	)
}
