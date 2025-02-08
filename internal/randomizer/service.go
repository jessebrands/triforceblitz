package randomizer

import (
	"errors"
	"io/fs"
	"log/slog"
	"path/filepath"
)

// Service provides methods to locate and load generators.
type Service struct {
	// Path where generators are installed.
	path string

	// Cached generators.
	generators []*Generator
}

func NewService(path string) *Service {
	return &Service{path: path}
}

func (s *Service) Synchronize() error {
	s.generators = []*Generator{}
	return filepath.WalkDir(s.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != MetadataFilename {
			return nil
		}
		metadata, err := OpenMetadata(path)
		if err != nil {
			slog.Warn("Could not load generator metadata, skipping.",
				"filename", path,
				"error", err)
			return filepath.SkipDir
		}
		version, err := VersionFromString(metadata.Version)
		if err != nil {
			slog.Warn("Could not parse generator version from metadata, skipping.",
				"version", metadata.Version,
				"filename", path,
				"error", err)
			return filepath.SkipDir
		}
		s.generators = append(s.generators, &Generator{
			Version: version,
			Path:    filepath.Dir(path),
			Presets: metadata.Presets,
		})
		return filepath.SkipDir
	})
}

// GetGenerators returns a list of all available generators.
func (s *Service) GetGenerators() []*Generator {
	return s.generators
}

// GetGenerator seeks a specific generator by version or returns an
// error if it cannot be found.
func (s *Service) GetGenerator(version Version) (*Generator, error) {
	for _, r := range s.generators {
		if version.Equal(r.Version) {
			return r, nil
		}
	}
	return nil, errors.New("generator not found")
}
