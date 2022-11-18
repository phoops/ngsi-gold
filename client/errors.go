package client

import (
	"github.com/pkg/errors"
)

// Client configuration
type ErrInvalidClientConfig error

var ErrMissingURL ErrInvalidClientConfig = errors.New("invalid client configuration: missing URL")
var ErrNegativeTimeout ErrInvalidClientConfig = errors.New("invalid client configuration: negative HTTP timeout")
var ErrWrongCustomHeaderFormat ErrInvalidClientConfig = errors.New("invalid client configuration: key or value of a custom header is empty")

// NGSI-LD errors URIs
var ngsiLdErrAlreadyExist = "https://uri.etsi.org/ngsi-ld/errors/AlreadyExists"
var ngsiLdErrBadData = "https://uri.etsi.org/ngsi-ld/errors/BadRequestData"
var ngsiLdErrInvalidRequest = "https://uri.etsi.org/ngsi-ld/errors/InvalidRequest"

// Operations
type ErrNgsiLdOperation error

var ErrNgsiLdEntityExists ErrNgsiLdOperation = errors.New("Entity already exists")
var ErrNgsiBadData ErrNgsiLdOperation = errors.New("Bad request")
var ErrNgsiLdInvalidRequest ErrNgsiLdOperation = errors.New("Invalid JSON of the request")
var ErrNgsiMixedResponse ErrNgsiLdOperation = errors.New("Mixed errors")
