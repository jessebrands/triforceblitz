package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v68/github"
)

var (
	cacheDir = os.Getenv("TRIFORCEBLITZ_PACKAGE_CACHE_DIR")
	manager  = NewPackageManager(os.Getenv("TRIFORCEBLITZ_GENERATORS_DIR"))
)

// listPackages lists all generator packages available.
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

	if branches != nil {
		whitelist = strings.Split(*branches, ",")
	}

	installer := NewInstaller(manager)
	installer.CachePackages = !*noCache
	if _, err := installer.InstallAll(whitelist); err != nil {
		panic(err)
	}
}

func main() {
	// Initialize the package manager.
	client := github.NewClient(nil)
	manager.AddSource(NewGitHubSource(client, "Elagatua", "OoT-Randomizer", cacheDir))
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
}
