package randomizer

import (
	"encoding/json"
	"io"
	"os"
)

// Metadata describes the metadata of a Generator.
type Metadata struct {
	Version    string   `json:"version"`
	Prerelease bool     `json:"prerelease,omitempty"`
	Presets    []Preset `json:"presets"`
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
	for _, p := range m.Presets {
		if p.Id == DefaultPreset {
			return nil
		}
	}
	return ErrNoDefaultPreset
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
