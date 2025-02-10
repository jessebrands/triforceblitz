package manager_test

import (
	"errors"
	"github.com/jessebrands/triforceblitz/internal/manager"
	"path/filepath"
	"testing"
)

func TestLockFile(t *testing.T) {
	t.Run("can be created", func(t *testing.T) {
		name := filepath.Join(t.TempDir(), "triforceblitz.lock")
		lock := manager.NewLockFile(name)
		if lock == nil {
			t.Errorf("expected non-nil lock file")
		}
	})

	t.Run("can lock", func(t *testing.T) {
		name := filepath.Join(t.TempDir(), "triforceblitz.lock")
		lf := manager.NewLockFile(name)
		err := lf.Lock(func() {})
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("calls function", func(t *testing.T) {
		name := filepath.Join(t.TempDir(), "triforceblitz.lock")
		lock := manager.NewLockFile(name)

		called := false
		_ = lock.Lock(func() {
			called = true
		})

		if !called {
			t.Errorf("expected function to be called")
		}
	})

	t.Run("can not acquire lock when locked", func(t *testing.T) {
		name := filepath.Join(t.TempDir(), "triforceblitz.lock")
		lf := manager.NewLockFile(name)
		var internalErr error
		err := lf.Lock(func() {
			internalErr = lf.Lock(func() {
				// Do nothing
			})
		})
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
		if !errors.Is(internalErr, manager.ErrLockFileLocked) {
			t.Errorf("expected ErrLockFileLocked, got: %v", internalErr)
		}
	})
}
