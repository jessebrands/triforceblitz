package randomizer

import (
	"errors"
	"io/fs"
	"log/slog"
	"path/filepath"
)

// Service provides methods to locate and load Randomizers.
type Service struct {
	// Path where Randomizers are installed.
	path string

	// Cached randomizers.
	randomizers []*Randomizer
}

func NewService(path string) *Service {
	return &Service{path: path}
}

func (s *Service) Synchronize() error {
	s.randomizers = []*Randomizer{}
	return filepath.WalkDir(s.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != MetadataFilename {
			return nil
		}
		metadata, err := OpenMetadata(path)
		if err != nil {
			slog.Warn("Could not load randomizer metadata, skipping.",
				"filename", path,
				"error", err)
			return filepath.SkipDir
		}
		version, err := VersionFromString(metadata.Version)
		if err != nil {
			slog.Warn("Could not parse randomizer version from metadata, skipping.",
				"version", metadata.Version,
				"filename", path,
				"error", err)
			return filepath.SkipDir
		}
		s.randomizers = append(s.randomizers, &Randomizer{
			Version: version,
			Path:    path,
			Presets: metadata.Presets,
		})
		return filepath.SkipDir
	})
}

// GetRandomizers returns a list of all available randomizers.
func (s *Service) GetRandomizers() []*Randomizer {
	return s.randomizers
}

// GetRandomizer seeks a specific version or returns an error if it cannot
// be found.
func (s *Service) GetRandomizer(version Version) (*Randomizer, error) {
	for _, r := range s.randomizers {
		if version.Equal(r.Version) {
			return r, nil
		}
	}
	return nil, errors.New("randomizer not found")
}
