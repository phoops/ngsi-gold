package client

import (
	"github.com/pkg/errors"
)

type ErrInvalidClientConfig error

var ErrMissingURL ErrInvalidClientConfig = errors.New("invalid client configuration: missing URL")
var ErrNegativeTimeout ErrInvalidClientConfig = errors.New("invalid client configuration: negative HTTP timeout")
var ErrWrongCustomHeaderFormat ErrInvalidClientConfig = errors.New("invalid client configuration: key or value of a custom header is empty")

// NGSI-LD errors URIs
var ngsiLdErrAlreadyExist = "https://uri.etsi.org/ngsi-ld/errors/AlreadyExists"
