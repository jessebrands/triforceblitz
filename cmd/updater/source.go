package main

import (
	"context"

	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

type Source interface {
	// Update updates the Package index.
	Update(context context.Context) error

	// GetAllPackages returns a list of available Packages from this Source.
	GetAllPackages() []Package

	// GetPackage gets a Package from the Source.
	GetPackage(version randomizer.Version) (Package, error)

	// DownloadPackage downloads a package from the Source to the cache.
	DownloadPackage(ctx context.Context, version randomizer.Version) error

	// UnpackPackage unpacks a package with the given version to the destination folder.
	UnpackPackage(ctx context.Context, version randomizer.Version, destination string) error

	// PurgePackage purges a package from the cache.
	PurgePackage(ctx context.Context, version randomizer.Version) error

	// IsCached returns whether the version is in the cache.
	IsCached(version randomizer.Version) bool

	// Type returns a string identifying the type of the Source.
	Type() string

	// Name returns a name uniquely identifying this Source.
	Name() string
}

func SourceIdentifier(source Source) string {
	return source.Type() + ":" + source.Name()
}
