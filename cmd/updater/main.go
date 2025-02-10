package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/jessebrands/triforceblitz/internal/config"
	"github.com/jessebrands/triforceblitz/internal/pkgman"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jessebrands/triforceblitz/internal/randomizer"

	"github.com/google/go-github/v68/github"
)

var (
	manager = pkgman.New()
)

// listPackages lists all generators packages available.
func listPackages() {
	packages := manager.AvailablePackages()

	// List the available packages:
	fmt.Printf("%-25s %-37s %-20s\n", "Version", "Installed", "Published at")
	for _, info := range packages {
		fmt.Printf("%-25.25s %-37.37v %-20.20s\n",
			info.Version.String(),
			info.IsInstalled(),
			info.PublishedAt.Format(time.RFC3339))
	}
}

func install(args []string) {
	var whitelist []string

	installFlags := flag.NewFlagSet("install", flag.ExitOnError)
	branches := installFlags.String("b", "blitz", "comma-separated list of branches to include.")
	noCache := installFlags.Bool("no-cache", false, "disable caching of package files.")
	if err := installFlags.Parse(args); err != nil {
		panic(err)
	}
	candidates := installFlags.Args()

	if branches != nil {
		whitelist = strings.Split(*branches, ",")
	}

	installer := pkgman.NewInstaller(manager)
	installer.CachePackages = !*noCache
	if len(candidates) > 0 {
		if _, err := installCandidates(installer, candidates); err != nil {
			fmt.Printf("Installation failed: %s\n", err.Error())
		}
	} else {
		if _, err := installer.InstallAll(whitelist); err != nil {
			fmt.Printf("Installation failed: %s\n", err.Error())
		}
	}
}

func installCandidates(installer *pkgman.Installer, candidates []string) ([]randomizer.Version, error) {
	// Turn candidates into versions:
	var versions []randomizer.Version
	for _, c := range candidates {
		version, err := randomizer.VersionFromString(c)
		if err != nil {
			fmt.Printf("Cannot select candidate %s, not a valid Triforce Blitz version: %s\n", c, err.Error())
			return nil, err
		}
		versions = append(versions, version)
	}
	return installer.Install(versions)
}

func main() {
	// Check that the output directories even exist.
	if err := os.MkdirAll(config.GetGeneratorDir(), 0755); err != nil {
		fmt.Printf("Cannot create generator directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(config.GetPackageCacheDir(), 0755); err != nil {
		fmt.Printf("Cannot create package cache directory: %v\n", err)
		os.Exit(1)
	}
	// Acquire a lockfile lock first.
	lockfile := pkgman.NewLockFile(config.GetLockFilename())
	err := lockfile.Lock(func() {
		// Initialize the package manager.
		client := github.NewClient(nil)
		manager.AddSource(pkgman.NewGitHubSource(client, "OoT-Randomizer", "Elagatua"))
		if err := manager.Update(context.Background()); err != nil {
			slog.Error("Could not refresh package index.", "error", err)
		}

		// Invoke the package manager.
		command := os.Args[1]
		switch command {
		case "list":
			listPackages()

		case "install":
			install(os.Args[2:])

		default:
			// Print out a useful help guide.
		}
	})
	if err != nil {
		if errors.Is(err, pkgman.ErrLockFileLocked) {
			fmt.Printf("Packages are locked, is another instance running?\n")
		} else if errors.Is(err, pkgman.ErrLockNotAcquired) {
			fmt.Printf("Could not acquire lock, do you have the right permissions?\n")
		}
	}
}
