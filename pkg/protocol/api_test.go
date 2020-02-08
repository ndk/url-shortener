package protocol

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortLinkRequest(t *testing.T) {
	Convey("Test validation", t, func() {
		r := CreateShortLinkRequest{}

		Convey("It fails if the URL is empty", func() {
			err := r.Bind(nil)
			assert.EqualError(t, err, "The URL mustn't be empty")
		})

		Convey("It fails if the URL is incorrect", func() {
			r.URL = "htt ttps://amazon.com"
			err := r.Bind(nil)
			assert.EqualError(t, err, "parse htt ttps://amazon.com: invalid URI for request")
		})

		Convey("It doesn't return any errors if everything is fine", func() {
			r.URL = "https://amazon.com"
			err := r.Bind(nil)
			assert.NoError(t, err)
		})
	})
}
