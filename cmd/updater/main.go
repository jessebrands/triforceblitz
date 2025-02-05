package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v68/github"
	"log"
	"os"
	"time"
)

var (
	manager = NewPackageManager(os.Getenv("TRIFORCEBLITZ_GENERATORS_DIR"))
)

// listPackages lists all generator packages available.
func listPackages() {
	packages, err := manager.AvailablePackages(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// List the available packages:
	fmt.Printf("%-25s %-37s %-20s\n", "Version", "Installed", "Published at")
	for version, info := range packages {
		fmt.Printf("%-25.25s %-37.37v %-20.20s\n",
			version.String(),
			info.Installed,
			info.PublishedAt.Format(time.RFC3339))
	}
}

func main() {
	// Initialize the package manager.
	client := github.NewClient(nil)
	manager.AddSource(NewGitHubSource(client, "Elagatua", "OoT-Randomizer"))

	// Invoke the package manager.
	command := os.Args[1]
	switch command {
	case "list":
		listPackages()

	default:
		// Print out a useful help guide.
	}
}
