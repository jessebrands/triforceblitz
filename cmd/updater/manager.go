package main

import (
	"context"
	"errors"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
)

type PackageIndex = map[generator.Version]PackageInfo

type PackageManager struct {
	installDir string
	sources    []Source
	index      PackageIndex
}

func NewPackageManager(installDir string) *PackageManager {
	return &PackageManager{
		installDir: installDir,
		index:      make(PackageIndex),
	}
}

func (m *PackageManager) AddSource(source Source) {
	m.sources = append(m.sources, source)
}

// Update updates the PackageManager's package index.
func (m *PackageManager) Update(ctx context.Context) error {
	for _, s := range m.sources {
		if err := s.Update(ctx); err != nil {
			slog.Warn("Could not update source package index. Skipping.",
				"source", SourceIdentifier(s),
				"error", err)
			continue
		}
		packages := s.GetAllPackages()
		for _, pkg := range packages {
			version := pkg.GetVersion()
			if info, ok := m.index[version]; !ok {
				m.index[version] = PackageInfo{
					Version:     version,
					PublishedAt: pkg.GetPublishedAt(),
					Sources:     []Source{s},
					Installed:   m.IsInstalled(version),
				}
			} else {
				info.Sources = append(info.Sources, s)
				if pkg.GetPublishedAt().Before(info.PublishedAt) {
					info.PublishedAt = pkg.GetPublishedAt()
				}
			}
		}
	}
	return nil
}

func (m *PackageManager) AvailablePackages() []PackageInfo {
	var packages []PackageInfo
	for _, info := range m.index {
		packages = append(packages, info)
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[j].PublishedAt.Before(packages[i].PublishedAt)
	})
	return packages
}

func (m *PackageManager) HasPackage(version generator.Version) bool {
	_, ok := m.index[version]
	return ok
}

// GetPackageDir returns the installation path for a specific package.
func (m *PackageManager) GetPackageDir(version generator.Version) string {
	return filepath.Join(m.installDir, version.String())
}

func (m *PackageManager) findEntrypoint(name string) (string, error) {
	entrypoint := ""
	found := false
	err := filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == generator.EntrypointFilename {
			entrypoint = path
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return entrypoint, err
	}
	if !found {
		return entrypoint, errors.New("entrypoint not found")
	}
	return entrypoint, nil
}

func (m *PackageManager) installPackageFromSource(ctx context.Context, s Source, version generator.Version) error {
	destination, err := os.MkdirTemp("", "TriforceBlitz")
	defer os.RemoveAll(destination)
	if err != nil {
		return err
	}
	if err := s.UnpackPackage(ctx, version, destination); err != nil {
		slog.Warn("Failed to retrieve package from source.",
			"version", version.String(),
			"source", SourceIdentifier(s),
			"error", err)
		return err
	}
	entrypoint, err := m.findEntrypoint(destination)
	if err != nil {
		slog.Warn("Failed to find entrypoint.",
			"destination", destination,
			"version", version.String(),
			"source", SourceIdentifier(s),
			"error", err)
		return err
	}
	generatorRoot := filepath.Dir(entrypoint)
	installDir := m.GetPackageDir(version)
	if err := os.CopyFS(installDir, os.DirFS(generatorRoot)); err != nil {
		return err
	}
	slog.Info("Successfully installed package.",
		"version", version.String(),
		"path", installDir)

	metadataFilename := filepath.Join(installDir, generator.MetadataFilename)
	slog.Info("Creating generator metadata for legacy versions.",
		"version", version.String(),
		"file", metadataFilename)
	return WriteMetadataFile(metadataFilename, version)
}

// Install attempts to install a generator.Generator to the installation directory managed
// by the PackageManager.
func (m *PackageManager) Install(ctx context.Context, version generator.Version) error {
	if m.IsInstalled(version) {
		return nil
	}
	info, ok := m.index[version]
	if !ok {
		return errors.New("package not available")
	}
	for _, s := range info.Sources {
		if err := m.installPackageFromSource(ctx, s, version); err != nil {
			continue
		}
		return nil
	}
	return errors.New("failed to install")
}

func (m *PackageManager) IsInstalled(version generator.Version) bool {
	path := m.GetPackageDir(version)
	metadata := filepath.Join(path, generator.MetadataFilename)
	entrypoint := filepath.Join(path, generator.EntrypointFilename)
	if _, err := os.Stat(metadata); err != nil {
		return false
	}
	if _, err := os.Stat(entrypoint); err != nil {
		return false
	}
	return true
}
