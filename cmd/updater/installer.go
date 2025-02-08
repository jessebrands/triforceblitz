package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

type Whitelist []string

func (w *Whitelist) Includes(s string) bool {
	return slices.Contains(*w, s)
}

// Installer provides a convenient interface over PackageManager for
// installing generator packages.
type Installer struct {
	manager       *PackageManager
	CachePackages bool
}

// NewInstaller creates a new Installer.
func NewInstaller(manager *PackageManager) *Installer {
	return &Installer{
		manager:       manager,
		CachePackages: true,
	}
}

// Install installs packages based on the given candidates. Returns a list
// of packages that were installed, if successful.
func (i *Installer) Install(versions []randomizer.Version) ([]randomizer.Version, error) {
	candidates, err := i.collect(versions)
	if err != nil {
		return []randomizer.Version{}, err
	}
	i.parallelInstall(context.Background(), candidates)
	return candidates, nil
}

// InstallAll installs all packages in the repository, filtered by a given
// list of branches to include. If the list is empty, all branches are
// included. Returns a list of packages that were installed.
func (i *Installer) InstallAll(whitelist Whitelist) ([]randomizer.Version, error) {
	candidates, err := i.collectAll(whitelist)
	if err != nil {
		return []randomizer.Version{}, err
	}
	return i.Install(candidates)
}

// collect takes a list of package versions and filters it down to packages
// that have not yet been installed.
//
// If the list of versions contains a nonexistent package, collect returns
// ErrPackageNotFound.
func (i *Installer) collect(versions []randomizer.Version) ([]randomizer.Version, error) {
	var candidates []randomizer.Version
	for _, v := range versions {
		if pkg, err := i.manager.GetPackage(v); err != nil {
			return []randomizer.Version{}, err
		} else if pkg.IsInstalled() {
			fmt.Printf("Generator %s is already installed\n", v.String())
			continue
		}
		candidates = append(candidates, v)
	}
	return candidates, nil
}

func (i *Installer) collectAll(whitelist Whitelist) ([]randomizer.Version, error) {
	var candidates []randomizer.Version
	for _, pkg := range i.manager.AvailablePackages() {
		if pkg.IsInstalled() || !whitelist.Includes(pkg.Version.Branch) {
			continue
		}
		candidates = append(candidates, pkg.Version)
	}
	return candidates, nil
}

func (i *Installer) install(ctx context.Context, version randomizer.Version) error {
	fmt.Printf("Selecting generator %s\n", version.String())
	if !i.manager.IsCached(version) {
		fmt.Printf("Downloading package %s\n", version.String())
		err := i.manager.Download(ctx, version)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Using cached package %s\n", version.String())
	}
	// Create a temporary directory to store the unpacked files.
	tempDir, err := os.MkdirTemp("", "TriforceBlitz")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	fmt.Printf("Unpacking package %s\n", version.String())
	if err := i.manager.Unpack(ctx, version, tempDir); err != nil {
		return err
	}
	if !i.CachePackages {
		fmt.Printf("Removing package %s from cache\n", version.String())
		if err := i.manager.Purge(ctx, version); err != nil {
			// This is not a fatal error, just inform the user and continue.
			fmt.Printf("Error removing package %s from cache: %s\n", version.String(), err.Error())
		}
	}
	fmt.Printf("Installing generator %s\n", version.String())
	if err := i.manager.Install(version, tempDir); err != nil {
		return err
	}
	fmt.Printf("Configuring generator %s\n", version.String())
	if err := i.manager.Configure(version); err != nil {
		return err
	}
	fmt.Printf("Installed generator %s\n", version.String())
	return nil
}

// parallelInstall takes a list of versions and installs them in parallel.
func (i *Installer) parallelInstall(ctx context.Context, versions []randomizer.Version) {
	var wg sync.WaitGroup
	for _, v := range versions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := i.install(ctx, v); err != nil {
				fmt.Printf("Failed to install generator %s: %s\n", v.String(), err.Error())
			}
		}()
	}
	wg.Wait()
}
