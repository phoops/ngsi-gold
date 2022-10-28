package client

import "time"

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
