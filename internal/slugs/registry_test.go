package slugs

import (
	"context"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
		s := &mockSlugifier{m: m}

		Convey("It fails if the slugifier has failed", func() {
			r := registry{
				slugifier:     s,
				instanceIndex: 5,
				slugsCount:    19,
			}

			m.On("NewSlug", int64(5), int64(19)).Return("", errors.New("NewSlug error"))

			_, err := r.RegisterURL(nil, "")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "NewSlug error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It fails if the value cannot be saved", func() {
			r := registry{
				slugifier:     s,
				instanceIndex: 5,
				slugsCount:    19,
				saveValue: func(key string, value string) error {
					args := m.Called(key, value)
					return args.Error(0)
				},
			}

			m.
				On("NewSlug", int64(5), int64(19)).Return("qwe", nil).
				On("1", "5:19", "http://en.wikipedia.com").Return(errors.New("saveValue error"))

			_, err := r.RegisterURL(context.TODO(), "http://en.wikipedia.com")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "saveValue error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It returns a new slug", func() {
			r := registry{
				slugifier:     s,
				instanceIndex: 5,
				slugsCount:    19,
				saveValue: func(key string, value string) error {
					args := m.Called(key, value)
					return args.Error(0)
				},
			}

			m.
				On("NewSlug", int64(5), int64(19)).Return("qwe", nil).
				On("1", "5:19", "http://en.wikipedia.com").Return(nil).
				On("NewSlug", int64(5), int64(20)).Return("asd", nil).
				On("1", "5:20", "http://en.wikipedia.com").Return(nil)

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
		s := &mockSlugifier{m: m}

		Convey("It fails if the slugifier has failed", func() {
			r := registry{
				slugifier:     s,
				instanceIndex: 5,
				slugsCount:    19,
				saveValue: func(key string, value string) error {
					args := m.Called(key, value)
					return args.Error(0)
				},
			}

			m.
				On("DecodeSlug", "123").Return(int64(321), int64(432), errors.New("DecodeSlug error"))

			_, err := r.GetURL(context.TODO(), "123")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "DecodeSlug error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It fails if the value cannot be loaded", func() {
			r := registry{
				slugifier:     s,
				instanceIndex: 5,
				slugsCount:    19,
				loadValue: func(key string) (string, error) {
					args := m.Called(key)
					return args.String(0), args.Error(1)
				},
			}

			m.
				On("DecodeSlug", "123").Return(int64(321), int64(432), nil).
				On("1", "321:432").Return("", errors.New("loadValue error"))

			_, err := r.GetURL(context.TODO(), "123")

			m.AssertExpectations(t)
			assert.EqualError(t, err, "loadValue error")
			assert.Equal(t, int64(5), r.instanceIndex)
			assert.Equal(t, int64(19), r.slugsCount)
		})

		Convey("It returns the correct URL", func() {
			r := registry{
				slugifier:     s,
				instanceIndex: 5,
				slugsCount:    19,
				loadValue: func(key string) (string, error) {
					args := m.Called(key)
					return args.String(0), args.Error(1)
				},
			}

			m.
				On("DecodeSlug", "123").Return(int64(321), int64(432), nil).
				On("1", "321:432").Return("http://uber.com", nil)

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
	s := &mockSlugifier{}
	r := NewRegistry(s, nil, nil, 178)
	assert.Equal(t,
		&registry{
			slugifier:     s,
			instanceIndex: 178,
		},
		r,
	)
}
