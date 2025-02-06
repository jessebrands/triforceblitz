package main

import "errors"

// ErrPackageNotFound is returned when a requested package does not exist.
var ErrPackageNotFound = errors.New("package not found")

// ErrDownloadFailed is returned when a package failed to download from all sources.
var ErrDownloadFailed = errors.New("download failed")

// ErrUnpackFailed is returned when a package fails to unpack.
var ErrUnpackFailed = errors.New("unpack failed")
