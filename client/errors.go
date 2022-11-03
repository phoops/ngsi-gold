package client

import (
	"github.com/pkg/errors"
)

type ErrInvalidClientConfig error

var ErrMissingURL ErrInvalidClientConfig = errors.New("invalid client configuration: missing URL")
var ErrNegativeTimeout ErrInvalidClientConfig = errors.New("invalid client configuration: negative HTTP timeout")
