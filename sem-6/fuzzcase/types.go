package fuzzcase

import "errors"

// ErrPort reports text that does not encode an allowed decimal port number.
var ErrPort = errors.New("invalid decimal port")
