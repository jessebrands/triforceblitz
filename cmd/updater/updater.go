package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v68/github"
	"log"
	"os"
	"time"
)

// listPackages lists all generator packages available.
func listPackages() {
	client := github.NewClient(nil)
	src := NewGitHubSource(client, "Elagatua", "OoT-Randomizer")
	packages, err := src.ListAvailable(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// List the available packages:
	fmt.Printf("%-25s %-37s %-20s\n", "Version", "Source", "Published at")
	for _, pkg := range packages {
		fmt.Printf("%-25.25s %-37.37s %-20.20s\n",
			pkg.Version,
			src.Type()+":"+src.Name(),
			pkg.PublishedAt.Format(time.RFC3339))
	}
}

func main() {
	command := os.Args[1]

	switch command {
	case "list":
		listPackages()

	default:
		// Print out a useful help guide.
	}
}
