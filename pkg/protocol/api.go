package protocol

import (
	"errors"
	"net/http"
	"net/url"
)

type ErrorResponse struct {
	Errors []Error `json:"errors,omitempty"`
}

type GeneralResponse struct {
	Data struct{} `json:"data,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////

type CreateShortLinkRequest struct {
	URL string `json:"url"`
}

func (c *CreateShortLinkRequest) Bind(r *http.Request) error {
	if c.URL == "" {
		return errors.New("The URL mustn't be empty")
	}
	if _, err := url.ParseRequestURI(c.URL); err != nil {
		return err
	}
	return nil
}

type CreateShortLinkResponse struct {
	Data struct {
		Slug string `json:"slug"`
	} `json:"data"`
}
