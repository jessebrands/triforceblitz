package main

import (
	"context"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"log/slog"
)

type PackageListing = map[generator.Version]PackageInfo

type PackageManager struct {
	installDir string
	sources    []Source
}

func NewPackageManager(installDir string) *PackageManager {
	return &PackageManager{installDir: installDir}
}

func (m *PackageManager) AddSource(source Source) {
	m.sources = append(m.sources, source)
}

func (m *PackageManager) AvailablePackages(ctx context.Context) (PackageListing, error) {
	listing := make(PackageListing)
	for _, s := range m.sources {
		packages, err := s.ListAvailable(ctx)
		if err != nil {
			slog.Warn("Could not get available packages from source. Skipping.",
				"source", SourceIdentifier(s),
				"error", err)
			continue
		}
		for _, pkg := range packages {
			if info, ok := listing[pkg.Version]; !ok {
				listing[pkg.Version] = PackageInfo{
					Version:     pkg.Version,
					PublishedAt: pkg.PublishedAt,
					Sources:     []Source{s},
					Installed:   false,
				}
			} else {
				info.Sources = append(info.Sources, s)
				if pkg.PublishedAt.Before(info.PublishedAt) {
					info.PublishedAt = pkg.PublishedAt
				}
			}
		}
	}
	return listing, nil
}
