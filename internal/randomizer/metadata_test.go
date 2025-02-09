package randomizer_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

const validMetadata = `
{
	"version": "7.23.1-blitz-1.59", 
	"presets": {
		"default": {
			"preset": "Triforce Blitz S2"
		},
		"triforce-blitz": {
			"preset": "Triforce Blitz",
			"ordinal": 100
		},
		"triforce-blitz-s2": {
			"preset": "Triforce Blitz S2",
			"ordinal": 200
		}
	}
}`

func TestUnmarshalMetadata(t *testing.T) {
	b := bytes.NewBufferString(validMetadata)
	metadata, err := randomizer.UnmarshalMetadata(b)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if metadata.Version != "7.23.1-blitz-1.59" {
		t.Errorf("expected version to be 7.23.1-blitz-1.59, got %s", metadata.Version)
	}
	if len(metadata.Presets) != 3 {
		t.Errorf("expected 3 presets, got %d", len(metadata.Presets))
	}
}

func TestMetadata_Validate(t *testing.T) {
	t.Run("must validate successfully", func(t *testing.T) {
		metadata := randomizer.Metadata{
			Version: "4.2.0-blitz-6.9",
			Presets: map[string]randomizer.Preset{
				"default":        {Value: "Triforce Blitz"},
				"triforce-blitz": {Value: "Triforce Blitz", Ordinal: 100},
			},
		}

		err := metadata.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("must error on invalid version", func(t *testing.T) {
		metadata := randomizer.Metadata{
			Version: "1.3-invalid-1",
			Presets: map[string]randomizer.Preset{
				"default": {Value: "Triforce Blitz"},
			},
		}

		err := metadata.Validate()
		if err == nil {
			t.Errorf("expected non-nil error")
		}
		if !errors.Is(err, randomizer.ErrInvalidVersion) {
			t.Errorf("expected ErrInvalidVersion error, got %v", err)
		}
	})

	t.Run("must error on missing default preset", func(t *testing.T) {
		metadata := randomizer.Metadata{
			Version: "1.0.0-blitz-1.0",
			Presets: map[string]randomizer.Preset{
				"triforce-blitz": randomizer.Preset{Value: "Triforce Blitz"},
			},
		}

		err := metadata.Validate()
		if err == nil {
			t.Errorf("expected non-nil error")
		}
		if !errors.Is(err, randomizer.ErrNoDefaultPreset) {
			t.Errorf("expected ErrNoDefaultPreset error, got %v", err)
		}
	})

	t.Run("must error when presets is empty", func(t *testing.T) {
		metadata := randomizer.Metadata{
			Version: "1.0.0-blitz-1.0",
			Presets: map[string]randomizer.Preset{},
		}

		err := metadata.Validate()
		if err == nil {
			t.Errorf("expected non-nil error")
		}
		if !errors.Is(err, randomizer.ErrNoPresets) {
			t.Errorf("expected ErrNoPresets error, got %v", err)
		}
	})
}
