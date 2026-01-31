package idgen

import "errors"

var (
	ErrInvalidNodeID = errors.New("node ID must be between 0 and 1023")
)
