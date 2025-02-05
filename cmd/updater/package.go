package main

import (
	"context"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"time"
)

type Package interface {
	// Download downloads the Package file to the destination.
	Download(ctx context.Context, destination string) error

	GetVersion() generator.Version
	GetPublishedAt() time.Time
}

type PackageInfo struct {
	Version     generator.Version
	PublishedAt time.Time
	Sources     []Source
	Installed   bool
}

func (info *PackageInfo) String() string {
	return info.Version.String()
}
