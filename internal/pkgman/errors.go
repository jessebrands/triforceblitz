package pkgman

import "errors"

// ErrLockFileLocked is returned when an attempt to Acquire a lock
// on a lock file failed because the lock file is already locked.
var ErrLockFileLocked = errors.New("lock already acquired")

// ErrLockNotAcquired is returned when an attempt to Acquire a lock
// file failed for any reason other than the lock file already being locked.
var ErrLockNotAcquired = errors.New("could not Acquire lock")

// ErrPackageNotFound is returned when a requested package does not exist.
var ErrPackageNotFound = errors.New("package not found")

// ErrDownloadFailed is returned when a package failed to download from all sources.
var ErrDownloadFailed = errors.New("download failed")

// ErrUnpackFailed is returned when a package fails to unpack.
var ErrUnpackFailed = errors.New("unpack failed")
