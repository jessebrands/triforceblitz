package generator_test

import (
	"errors"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"reflect"
	"testing"
)

func TestVersionFromString(t *testing.T) {
	t.Run("must parse successfully", func(t *testing.T) {
		version, err := generator.VersionFromString("4.2.0-blitz-6.9")
		want := generator.Version{Major: 4, Minor: 2, Branch: "blitz", BranchMajor: 6, BranchMinor: 9}
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(version, want) {
			t.Errorf("got %+v, want %+v", version, want)
		}
	})

	t.Run("must return ErrInvalidVersion when invalid", func(t *testing.T) {
		_, err := generator.VersionFromString("4.2-invalid-01")
		if err == nil {
			t.Errorf("expected non-nil error")
		}
		if !errors.Is(err, generator.ErrInvalidVersion) {
			t.Errorf("expected %v, got %v", generator.ErrInvalidVersion, err)
		}
	})
}
