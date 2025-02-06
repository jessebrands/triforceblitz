package main

import "errors"

// ErrPackageNotFound is returned when a requested package does not exist.
var ErrPackageNotFound = errors.New("package not found")
