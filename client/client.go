package client

import "net/http"

type NgsiLdClient struct {
	c   *http.Client
	url string
}

// OptionFunc is a function that configures a NgsiLdClient.
// Provide some of them to New to tweak client behaviour
type OptionFunc func(*NgsiLdClient) error

func New(options ...OptionFunc) (*NgsiLdClient, error) {
	ngsiLdClient := &NgsiLdClient{}

	ngsiLdClient.c = &http.Client{}

	// apply the options
	for _, option := range options {
		if err := option(ngsiLdClient); err != nil {
			return nil, err
		}
	}

	if err := ngsiLdClient.Validate(); err != nil {
		return nil, err
	}

	return ngsiLdClient, nil
}

func (c *NgsiLdClient) Validate() error {
	if c.url == "" {
		return ErrMissingURL
	}
	return nil
}
