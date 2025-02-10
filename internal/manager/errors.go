package manager

import "errors"

// ErrLockFileLocked is returned when an attempt to acquire a lock
// on a lock file failed because the lock file is already locked.
var ErrLockFileLocked = errors.New("lock already acquired")

// ErrLockNotAcquired is returned when an attempt to acquire a lock
// file failed for any reason other than the lock file already being locked.
var ErrLockNotAcquired = errors.New("could not acquire lock")
