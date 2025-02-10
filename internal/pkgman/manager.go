package pkgman

import (
	"context"
	"errors"
	"github.com/jessebrands/triforceblitz/internal/config"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

type PackageIndex = map[randomizer.Version]PackageInfo

type PackageManager struct {
	installDir string
	sources    []Source
	index      PackageIndex
}

func New() *PackageManager {
	return &PackageManager{
		installDir: config.GetGeneratorDir(),
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

func (m *PackageManager) GetPackage(version randomizer.Version) (PackageInfo, error) {
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

func (m *PackageManager) HasPackage(version randomizer.Version) bool {
	_, err := m.GetPackage(version)
	return err == nil
}

// installationDir returns the installation directory for a specific package.
func (m *PackageManager) installationDir(version randomizer.Version) string {
	return filepath.Join(m.installDir, version.String())
}

func (m *PackageManager) findEntrypoint(name string) (string, error) {
	entrypoint := ""
	found := false
	err := filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == randomizer.EntrypointFilename {
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

func (m *PackageManager) IsCached(version randomizer.Version) bool {
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
func (m *PackageManager) Download(ctx context.Context, version randomizer.Version) error {
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

func (m *PackageManager) Unpack(ctx context.Context, version randomizer.Version, destination string) error {
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

// Install copies the generator files from the source directory to the installation directory for the given version
// and updates the package information in the index.
func (m *PackageManager) Install(version randomizer.Version, sourceDir string) error {
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
func (m *PackageManager) Configure(version randomizer.Version) error {
	pkg, err := m.GetPackage(version)
	if err != nil {
		return err
	}
	// Older generators may not have the metadata file included, in that case
	// we will want to generate it for them based on some preset data.
	metadataFilename := filepath.Join(pkg.installDir, randomizer.MetadataFilename)
	if _, err := os.Stat(metadataFilename); os.IsNotExist(err) {
		if err := CreateMetadataFile(metadataFilename, version); err != nil {
			return err
		}
	}
	// The randomizer comes with a bunch of binaries.
	// Without making these executable, the program won't run.
	binaries := []string{
		filepath.Join(pkg.installDir, "Compress/Compress"),
		filepath.Join(pkg.installDir, "Compress/Compress.exe"),
		filepath.Join(pkg.installDir, "Decompress/Decompress"),
		filepath.Join(pkg.installDir, "Decompress/Decompress.exe"),
		filepath.Join(pkg.installDir, "bin/Compress/Compress"),
		filepath.Join(pkg.installDir, "bin/Compress/Compress.exe"),
		filepath.Join(pkg.installDir, "bin/Decompress/Decompress"),
		filepath.Join(pkg.installDir, "bin/Decompress/Decompress.exe"),
	}
	for _, name := range binaries {
		if err := os.Chmod(name, 0755); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return err
		}
	}
	return nil
}

// Purge removes a package from the cache.
func (m *PackageManager) Purge(ctx context.Context, version randomizer.Version) error {
	pkg, err := m.GetPackage(version)
	if err != nil {
		return err
	}
	for _, source := range pkg.Sources {
		if err := source.PurgePackage(ctx, version); err != nil {
			return err
		}
	}
	return nil
}
