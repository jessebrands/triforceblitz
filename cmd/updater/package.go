package main

import (
	"github.com/jessebrands/triforceblitz/internal/generator"
	"time"
)

type Package struct {
	Version     generator.Version
	PublishedAt time.Time
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
