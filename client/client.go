package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/phoops/ngsild/ldcontext"
)

var (
	ngsiLdUserAgent = requestHeader{key: "User-Agent", value: "ngsild-client"}
	jsonLdBody      = requestHeader{key: "Content-Type", value: "application/ld+json"}
	jsonResponse    = requestHeader{key: "Accept", value: "application/json"}
)

type NgsiLdClient struct {
	c       *http.Client
	url     string
	headers map[string]string
}

// OptionFunc is a function that configures a NgsiLdClient.
// Provide some of them to New to tweak client behaviour
type OptionFunc func(*NgsiLdClient) error

func New(options ...OptionFunc) (*NgsiLdClient, error) {
	ngsiLdClient := &NgsiLdClient{}

	ngsiLdClient.c = &http.Client{}
	ngsiLdClient.headers = map[string]string{}

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

type requestBody map[string]json.RawMessage
type requestError struct {
	ErrType string `json:"type"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
}

func addContext(payload json.Marshaler, ldCtx *ldcontext.LdContext) (requestBody, error) {
	requestBody := requestBody{}
	// struct -> json
	serializedPayload, err := payload.MarshalJSON()
	if err != nil {
		return nil, err
	}
	// json -> key-value
	err = json.Unmarshal(serializedPayload, &requestBody)
	if err != nil {
		return nil, err
	}

	// Add Context
	serializedContext, err := json.Marshal(ldCtx)
	if err != nil {
		return nil, err
	}
	requestBody["@context"] = serializedContext

	return requestBody, nil
}

type requestHeader struct {
	key   string
	value string
}

func (c *NgsiLdClient) newRequest(ctx context.Context, method, url string, body io.Reader, headers ...requestHeader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// expectation of a ngsi-ld client
	headers = append(headers, ngsiLdUserAgent)
	headers = append(headers, jsonResponse)

	// set the global headers
	for header, value := range c.headers {
		req.Header.Add(header, value)
	}

	for _, ah := range headers {
		req.Header.Add(ah.key, ah.value)
	}
	return req, nil
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

func SetGlobalHeader(key string, value string) OptionFunc {
	return func(client *NgsiLdClient) error {
		if key == "" || value == "" {
			return ErrWrongCustomHeaderFormat
		}
		client.headers[key] = value
		return nil
	}
}
