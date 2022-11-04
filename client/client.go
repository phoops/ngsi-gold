package client

import (
	"net/http"
	"time"
)

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

// SetURL makes the client connect to the specified Context Broker.
func SetURL(url string) OptionFunc {
	return func(c *NgsiLdClient) error {
		c.url = url
		return nil
	}
}

// SetClientTimeout specifies a value for HTTP client timeout
func SetClientTimeout(timeout time.Duration) OptionFunc {
	return func(client *NgsiLdClient) error {
		if timeout <= 0 {
			return ErrNegativeTimeout
		}
		client.c.Timeout = timeout
		return nil
	}
}
