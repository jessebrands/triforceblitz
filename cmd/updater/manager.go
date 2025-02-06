package main

import (
	"context"
	"errors"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"io/fs"
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
					installDir:  m.installationDir(version),
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

func (m *PackageManager) GetPackage(version generator.Version) (PackageInfo, error) {
	if info, ok := m.index[version]; ok {
		return info, nil
	} else {
		return info, ErrPackageNotFound
	}
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
	_, err := m.GetPackage(version)
	return err == nil
}

// installationDir returns the installation directory for a specific package.
func (m *PackageManager) installationDir(version generator.Version) string {
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

func (m *PackageManager) IsCached(version generator.Version) bool {
	info, ok := m.index[version]
	if !ok {
		return false
	}
	for _, s := range info.Sources {
		if s.IsCached(version) {
			return true
		}
	}
	return false
}

// Download downloads a Package to the cache directory.
func (m *PackageManager) Download(ctx context.Context, version generator.Version) error {
	info, ok := m.index[version]
	if !ok {
		return ErrPackageNotFound
	}
	for _, s := range info.Sources {
		err := s.DownloadPackage(ctx, version)
		if err != nil {
			continue
		}
		return nil
	}
	return ErrDownloadFailed
}

func (m *PackageManager) Unpack(ctx context.Context, version generator.Version, destination string) error {
	if pkg, err := m.GetPackage(version); err != nil {
		return err
	} else {
		for _, s := range pkg.Sources {
			if !s.IsCached(version) {
				continue
			}
			if err := s.UnpackPackage(ctx, version, destination); err != nil {
				continue
			}
			return nil
		}
		return ErrUnpackFailed
	}
}

// Install copies the Generator files from the source directory to the installation directory for the given version
// and updates the package information in the index.
func (m *PackageManager) Install(version generator.Version, sourceDir string) error {
	pkg, err := m.GetPackage(version)
	if err != nil {
		return err
	}
	entrypoint, err := m.findEntrypoint(sourceDir)
	if err != nil {
		return err
	}
	return os.CopyFS(pkg.installDir, os.DirFS(filepath.Dir(entrypoint)))
}

// Configure sets up a generator so that it can be used.
func (m *PackageManager) Configure(version generator.Version) error {
	pkg, err := m.GetPackage(version)
	if err != nil {
		return err
	}
	// Older generators may not have the metadata file included, in that case
	// we will want to generate it for them based on some preset data.
	metadataFilename := filepath.Join(pkg.installDir, generator.MetadataFilename)
	if _, err := os.Stat(metadataFilename); os.IsNotExist(err) {
		return CreateMetadataFile(metadataFilename, version)
	}
	return nil
}
