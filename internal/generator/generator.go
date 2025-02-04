package generator

import (
	"encoding/json"
	"io"
	"io/fs"
	"path/filepath"
)

const (
	EntrypointFilename = "OoTRandomizer.py"
	MetadataFilename   = "Generator.json"
)

type Generator struct {
	Path    string
	Version string
	Presets []Preset
}

type Preset struct {
	Id      string
	Preset  string
	Ordinal int
}

type Metadata struct {
	Version string
	Presets []Preset
}

func unmarshalMetadata(r io.Reader) (*Metadata, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	metadata := new(Metadata)
	if err := json.Unmarshal(b, metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func loadMetadata(fsys fs.FS, name string) (*Metadata, error) {
	file, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return unmarshalMetadata(file)
}

// FindGeneratorsFromFS recursively searches the filesystem for generators.
func FindGeneratorsFromFS(fsys fs.FS, root string) ([]*Generator, error) {
	var generators []*Generator
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if d.Name() != MetadataFilename {
			return nil
		}
		metadata, err := loadMetadata(fsys, path)
		if err != nil {
			return err
		}
		generators = append(generators, &Generator{
			Path:    filepath.Dir(path),
			Version: metadata.Version,
			Presets: metadata.Presets,
		})
		return nil
	})
	return generators, err
}

// Entrypoint returns the path to the Generator entrypoint file.
func (g *Generator) Entrypoint() string {
	return filepath.Join(g.Path, EntrypointFilename)
}
