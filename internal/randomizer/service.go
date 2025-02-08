package randomizer

import (
	"io/fs"
	"log/slog"
	"path/filepath"
)

// Service provides methods to locate and load Randomizers.
type Service struct {
	// Path where Randomizers are installed.
	path string
}

func NewService(path string) *Service {
	return &Service{path: path}
}

// GetRandomizers returns a list of all available randomizers.
func (s *Service) GetRandomizers() ([]*Randomizer, error) {
	var randomizers []*Randomizer
	err := filepath.WalkDir(s.path, func(path string, d fs.DirEntry, e error) error {
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
		randomizers = append(randomizers, &Randomizer{
			Version: version,
			Path:    path,
			Presets: metadata.Presets,
		})
		return filepath.SkipDir
	})
	if err != nil {
		return randomizers, err
	}
	return randomizers, nil
}
