package config

import (
	"os"
	"path/filepath"
)

const (
	// EnvGeneratorDir is the name of the environment variable that stores the
	// directory where generators should be installed.
	EnvGeneratorDir = "TRIFORCEBLITZ_GENERATORS_DIR"

	// EnvPackageCacheDir is the name of the environment variable that stores
	// the directory where packages should be stored to.
	EnvPackageCacheDir = "TRIFORCEBLITZ_PACKAGE_CACHE_DIR"

	// EnvLockFilename is the name of the environment variable that stores the
	// filename of the lockfile.
	EnvLockFilename = "TRIFORCEBLITZ_LOCK_FILE"
)

// GetGeneratorDir returns the directory where generators should be stored.
func GetGeneratorDir() string {
	if path := os.Getenv(EnvGeneratorDir); path != "" {
		return path
	}
	if path, err := os.UserCacheDir(); err == nil {
		return filepath.Join(path, "triforceblitz/generators")
	} else {
		return "generators"
	}
}

// GetPackageCacheDir returns the path to the directory where packages should be
// stored.
func GetPackageCacheDir() string {
	if path := os.Getenv(EnvPackageCacheDir); path != "" {
		return path
	}
	if path, err := os.UserCacheDir(); err == nil {
		return filepath.Join(path, "triforceblitz/packages")
	} else {
		return "packages"
	}
}

func GetLockFilename() string {
	if path := os.Getenv(EnvLockFilename); path != "" {
		return path
	}
	// If we don't have an environment variable, just use the package dir.
	return filepath.Join(GetPackageCacheDir(), "triforceblitz.lock")
}
