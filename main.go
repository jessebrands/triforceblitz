package main

import (
	"log"

	"github.com/jessebrands/triforceblitz/python"
	"github.com/jessebrands/triforceblitz/randomizer"
)

func main() {
	interpreter, err := python.Find()
	if err != nil {
		log.Fatalln("Could not find Python interpreter:", err)
	}
	version, err := interpreter.Version()
	if err != nil {
		log.Fatalln("Could not get Python version:", err)
	}
	log.Println("Found Python version", version, "at", interpreter.Path())
	// Initialize our downloader and check what versions are available.
	dl := randomizer.NewDownloader("Elagatua", "OoT-Randomizer", "/home/bee")
	releases, err := dl.GetAvailableReleases()
	if err != nil {
		log.Fatalln("Could not get Triforce Blitz releases:", err)
	}
	for i, r := range releases {
		log.Printf("%d. %s", i+1, r.Version)
	}
	tarball, err := dl.Download(releases[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Saved release tarball for version", releases[0].Version, "to", tarball)
	if err := dl.Install(tarball); err != nil {
		log.Fatalln("Failed to install randomizer:", err)
	}
}
