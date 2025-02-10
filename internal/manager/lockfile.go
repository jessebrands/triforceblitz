package manager

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
)

type LockFile struct {
	name string
}

// acquire attempts to acquire the lock file.
//
// # If the lock file already exists, it will return ErrLockFileLocked
//
// If the lock file cannot be created,
func (l *LockFile) acquire() error {
	if locked, err := l.Locked(); err != nil {
		return err
	} else if locked {
		return ErrLockFileLocked
	}
	f, err := os.OpenFile(l.name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return ErrLockNotAcquired
	}
	defer f.Close()
	pid := strconv.Itoa(os.Getpid())
	if _, err := f.Write([]byte(pid)); err != nil {
		return err
	}
	return nil
}

func (l *LockFile) release() error {
	return os.Remove(l.name)
}

func (l *LockFile) Lock(f func()) error {
	if err := l.acquire(); err != nil {
		return err
	}
	f()
	return l.release()
}

func (l *LockFile) Locked() (bool, error) {
	if _, err := os.Stat(l.name); err == nil {
		return true, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return true, err
	}
	return false, nil
}

func NewLockFile(name string) *LockFile {
	return &LockFile{name}
}
