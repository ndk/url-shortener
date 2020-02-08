package slugs

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/speps/go-hashids"
	"github.com/stretchr/testify/assert"
)

func TestHashidsSlugifier(t *testing.T) {
	Convey("Test HashidsSlugifier", t, func() {
		Convey("It is constractable", func() {
			cfg := &Config{}
			s, err := NewHashidsSlugifier(cfg)
			assert.NoError(t, err)
			assert.NotNil(t, s)
		})

		Convey("It returns decodable slug", func() {
			s, err := NewHashidsSlugifier(&Config{"123", 8})
			assert.NoError(t, err)
			slug, err := s.NewSlug(123, 456)
			assert.NoError(t, err)
			assert.Equal(t, slug, "2V87tRpL")
			instanceIndex, slugIndex, err := s.DecodeSlug(slug)
			assert.NoError(t, err)
			assert.Equal(t, int64(123), instanceIndex)
			assert.Equal(t, int64(456), slugIndex)
		})

		Convey("Creating of a new slug fails if hashids has failed", func() {
			s, err := NewHashidsSlugifier(&Config{"123", 8})
			assert.NoError(t, err)
			_, err = s.NewSlug(123, -456)
			assert.EqualError(t, err, "negative number not supported")
		})

		Convey("Test decoding", func() {
			s, err := NewHashidsSlugifier(&Config{"123", 8})
			assert.NoError(t, err)

			Convey("It fails if decoding has failed", func() {
				_, _, err := s.DecodeSlug("123")
				assert.EqualError(t, err, "mismatch between encode and decode: 123 start 4rlRNlnd re-encoded. result: [0]")
			})

			Convey("The count of the numbers must be equal to 2", func() {
				hd := hashids.NewData()
				hd.Salt = "123"
				hd.MinLength = 8
				h, err := hashids.NewWithData(hd)
				assert.NoError(t, err)
				slug, err := h.EncodeInt64([]int64{123, 456, 789})
				assert.NoError(t, err)
				_, _, err = s.DecodeSlug(slug)
				assert.EqualError(t, err, "The slug is corrupted")
			})
		})
	})
}
