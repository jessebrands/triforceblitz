package randomizer

import (
	"io/fs"
	"log/slog"
	"path/filepath"
)

// Service provides methods to locate and load Generators.
type Service struct {
	// Path where Generators are installed.
	path string
}

func NewService(path string) *Service {
	return &Service{path: path}
}

// GetGenerators returns a list of all available Generators.
func (s *Service) GetGenerators() ([]*Generator, error) {
	var generators []*Generator
	err := filepath.WalkDir(s.path, func(path string, d fs.DirEntry, err error) error {
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
		generators = append(generators, &Generator{
			Version: version,
			Path:    path,
			Presets: metadata.Presets,
		})
		return filepath.SkipDir
	})
	if err != nil {
		return generators, err
	}
	return generators, nil
}
