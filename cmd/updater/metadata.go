package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

var (
	// blitz-0.1 - blitz-0.19
	metadataPrerelease = randomizer.Metadata{
		Prerelease: true,
		Presets: randomizer.PresetMap{
			randomizer.DefaultPreset: {Value: "Triforce Blitz"},
		},
	}
	// blitz-0.20 - blitz-0.37
	metadataSeason1 = randomizer.Metadata{
		Prerelease: false,
		Presets: randomizer.PresetMap{
			randomizer.DefaultPreset: {Value: "Triforce Blitz"},
			"triforce-blitz-s1":      {Value: "Triforce Blitz", Ordinal: 100},
		},
	}
	// blitz-0.40 - blitz-0.42
	metadataSeason2 = randomizer.Metadata{
		Prerelease: false,
		Presets: randomizer.PresetMap{
			randomizer.DefaultPreset: {Value: "Triforce Blitz S2"},
			"triforce-blitz-s2":      {Value: "Triforce Blitz S2", Ordinal: 200},
		},
	}
	// blitz-0.43
	metadataSeason2Post = randomizer.Metadata{
		Prerelease: false,
		Presets: randomizer.PresetMap{
			randomizer.DefaultPreset:          {Value: "Triforce Blitz S2"},
			"triforce-blitz-s2-jabus-revenge": {Value: "Triforce Blitz S2", Ordinal: 201},
		},
	}
	// blitz-0.50 - blitz-0.59
	metadataSeason3 = randomizer.Metadata{
		Prerelease: false,
		Presets: randomizer.PresetMap{
			randomizer.DefaultPreset: {Value: "Triforce Blitz S3"},
			"triforce-blitz-s3":      {Value: "Triforce Blitz S3", Ordinal: 300},
		},
	}
)

func injectVersionIntoMetadata(version randomizer.Version, metadata randomizer.Metadata) randomizer.Metadata {
	metadata.Version = version.String()
	return metadata
}

// GetMetadataForVersion returns Metadata for the Triforce Blitz version requested. The function
// will return an error in case no Metadata is available for that version. Newer versions of
// Triforce Blitz should ship with their own metadata files.
func GetMetadataForVersion(version randomizer.Version) (randomizer.Metadata, error) {
	if version.Branch != "blitz" {
		return randomizer.Metadata{}, errors.New("unsupported branch")
	}
	if version.BranchMajor > 0 || version.BranchMinor > 59 {
		return randomizer.Metadata{}, errors.New("only versions before blitz-0.60 are supported")
	}
	if version.BranchMinor < 20 {
		return injectVersionIntoMetadata(version, metadataPrerelease), nil
	}
	if version.BranchMinor < 40 {
		return injectVersionIntoMetadata(version, metadataSeason1), nil
	}
	if version.BranchMinor < 43 {
		return injectVersionIntoMetadata(version, metadataSeason2), nil
	}
	if version.BranchMinor == 43 {
		return injectVersionIntoMetadata(version, metadataSeason2Post), nil
	}
	if version.BranchMinor < 60 {
		return injectVersionIntoMetadata(version, metadataSeason3), nil
	}
	return randomizer.Metadata{}, fmt.Errorf("unsupported version: %s", version.String())
}

func CreateMetadataFile(name string, version randomizer.Version) error {
	metadata, err := GetMetadataForVersion(version)
	if err != nil {
		// Unsupported version, not an error.
		return nil
	}
	b, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return err
	}
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}
