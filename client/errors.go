package client

import (
	"github.com/pkg/errors"
)

var ErrMissingURL = errors.New("invalid client configuration: missing URL")
var ErrNegativeTimeout = errors.New("invalid client configuration: negative HTTP timeout")
