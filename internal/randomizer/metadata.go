package randomizer

import (
	"encoding/json"
	"io"
	"os"
)

const (
	MetadataFilename = ".generator-metadata.json"
)

// Metadata describes details about a Generator.
type Metadata struct {
	Version    string    `json:"version"`
	Prerelease bool      `json:"prerelease,omitempty"`
	Presets    PresetMap `json:"presets"`
}

// Validate validates the contents of the Metadata.
//
// If the Metadata contains an invalid version, Validate will return an
// ErrInvalidVersion error.
//
// If the Metadata is missing the mandatory 'default' Preset, Validate
// will return an ErrNoDefaultPreset error.
//
// If the Metadata contains no presets at all, Validate will return an
// ErrNoPresets error.
func (m *Metadata) Validate() error {
	if !ValidVersion(m.Version) {
		return ErrInvalidVersion
	}
	if len(m.Presets) == 0 {
		return ErrNoPresets
	}
	if _, err := m.Presets.Default(); err != nil {
		return err
	}
	return nil
}

// UnmarshalMetadata takes an io.Reader containing a Metadata file and parses it.
//
// This function performs validation on the Metadata by calling Metadata.Validate. If
// validation fails for the Metadata, then UnmarshalMetadata will return one of the
// errors returned by Metadata.Validate.
func UnmarshalMetadata(r io.Reader) (*Metadata, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	metadata := new(Metadata)
	if err := json.Unmarshal(b, metadata); err != nil {
		return nil, err
	}
	return metadata, metadata.Validate()
}

// OpenMetadata opens a Metadata file and parses it.
func OpenMetadata(name string) (*Metadata, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return UnmarshalMetadata(f)
}
