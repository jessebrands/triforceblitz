package pkgman

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

type Package interface {
	// Download downloads the Package file to the destination.
	Download(ctx context.Context, destination string) error
	GetVersion() randomizer.Version
	GetPublishedAt() time.Time
}

type PackageInfo struct {
	Version     randomizer.Version
	PublishedAt time.Time
	Sources     []Source
	installDir  string
}

func (info *PackageInfo) IsInstalled() bool {
	metadata := filepath.Join(info.installDir, randomizer.MetadataFilename)
	entrypoint := filepath.Join(info.installDir, randomizer.EntrypointFilename)
	if _, err := os.Stat(metadata); err != nil {
		return false
	}
	if _, err := os.Stat(entrypoint); err != nil {
		return false
	}
	return true
}

func (info *PackageInfo) GetInstallDir() string {
	return info.installDir
}

func (info *PackageInfo) String() string {
	return info.Version.String()
}
