package main

import (
	"context"
	"fmt"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"slices"
	"sync"
)

type Whitelist []string

func (w *Whitelist) Includes(s string) bool {
	return len(*w) >= 1 && slices.Contains(*w, s)
}

// Installer provides a convenient interface over PackageManager for
// installing generator packages.
type Installer struct {
	manager *PackageManager
}

// NewInstaller creates a new Installer.
func NewInstaller(manager *PackageManager) *Installer {
	return &Installer{
		manager: manager,
	}
}

// Install installs packages based on the given candidates. Returns a list
// of packages that were installed, if successful.
func (i *Installer) Install(versions []generator.Version) ([]generator.Version, error) {
	candidates, err := i.collect(versions)
	if err != nil {
		return []generator.Version{}, err
	}
	i.parallelInstall(context.Background(), candidates)
	return candidates, nil
}

// InstallAll installs all packages in the repository, filtered by a given
// list of branches to include. If the list is empty, all branches are
// included. Returns a list of packages that were installed.
func (i *Installer) InstallAll(whitelist Whitelist) ([]generator.Version, error) {
	candidates, err := i.collectAll(whitelist)
	if err != nil {
		return []generator.Version{}, err
	}
	return i.Install(candidates)
}

// collect takes a list of package versions and filters it down to packages
// that have not yet been installed.
//
// If the list of versions contains a nonexistent package, collect returns
// ErrPackageNotFound.
func (i *Installer) collect(versions []generator.Version) ([]generator.Version, error) {
	var candidates []generator.Version
	for _, v := range versions {
		if i.manager.IsInstalled(v) {
			fmt.Printf("Package '%s' is already installed\n", v.String())
			continue
		}
		if !i.manager.HasPackage(v) {
			return candidates, ErrPackageNotFound
		}
		candidates = append(candidates, v)
	}
	return candidates, nil
}

func (i *Installer) collectAll(whitelist Whitelist) ([]generator.Version, error) {
	var candidates []generator.Version
	for _, pkg := range i.manager.AvailablePackages() {
		if pkg.Installed || !whitelist.Includes(pkg.Version.Branch) {
			continue
		}
		candidates = append(candidates, pkg.Version)
	}
	return candidates, nil
}

// parallelInstall takes a list of versions and installs them in parallel.
func (i *Installer) parallelInstall(ctx context.Context, versions []generator.Version) {
	var wg sync.WaitGroup
	for _, v := range versions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := i.manager.Install(ctx, v); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}
