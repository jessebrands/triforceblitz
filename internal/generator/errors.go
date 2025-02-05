package generator

import "errors"

// ErrInvalidVersion is returned when a Version string does not match the
// expected format for a Version.
var ErrInvalidVersion = errors.New("invalid version")
