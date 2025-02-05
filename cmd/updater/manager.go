package main

import (
	"context"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"log/slog"
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
					Installed:   false,
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
	destination := os.TempDir()
	for _, s := range info.Sources {
		if err := s.UnpackPackage(ctx, version, destination); err != nil {
			continue
		}
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
