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
		packages, err := s.ListAvailable(ctx)
		if err != nil {
			slog.Warn("Could not get available packages from source. Skipping.",
				"source", SourceIdentifier(s),
				"error", err)
			continue
		}
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
