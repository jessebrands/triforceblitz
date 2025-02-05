package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v68/github"
	"log/slog"
	"os"
	"time"
)

var (
	manager = NewPackageManager(os.Getenv("TRIFORCEBLITZ_GENERATORS_DIR"))
)

// listPackages lists all generator packages available.
func listPackages() {
	packages := manager.AvailablePackages()

	// List the available packages:
	fmt.Printf("%-25s %-37s %-20s\n", "Version", "Installed", "Published at")
	for _, info := range packages {
		fmt.Printf("%-25.25s %-37.37v %-20.20s\n",
			info.Version.String(),
			info.Installed,
			info.PublishedAt.Format(time.RFC3339))
	}
}

func main() {
	// Initialize the package manager.
	client := github.NewClient(nil)
	manager.AddSource(NewGitHubSource(client, "Elagatua", "OoT-Randomizer"))
	if err := manager.Update(context.Background()); err != nil {
		slog.Error("Could not refresh package index.", "error", err)
	}

	// Invoke the package manager.
	command := os.Args[1]
	switch command {
	case "list":
		listPackages()

	default:
		// Print out a useful help guide.
	}
}
