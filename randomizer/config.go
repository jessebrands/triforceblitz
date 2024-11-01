package randomizer

import "os"

const (
	envRepositoryOwner = "TRIFORCE_BLITZ_REPOSITORY_OWNER"
	envRepositoryName  = "TRIFORCE_BLITZ_REPOSITORY_NAME"
	envDownloadPath    = "TRIFORCE_BLITZ_DOWNLOAD_PATH"
	envInstallPath     = "TRIFORCE_BLITZ_INSTALL_PATH"
)

type repository struct {
	owner string
	name  string
}

// Config describes the parameters of the randomizer service.
type Config struct {
	repository   repository
	downloadPath string
	installPath  string
}

// envOr gets an environment variable, or a default value if not found.
func envOr(key string, value string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return value
}

// DefaultConfig returns a configuration created from environment variables, if possible, and
// otherwise returns sensible defaults instead.
func DefaultConfig() *Config {
	return &Config{
		repository: repository{
			owner: envOr(envRepositoryOwner, "Elagatua"),
			name:  envOr(envRepositoryName, "OoT-Randomizer"),
		},
		downloadPath: envOr(envDownloadPath, "~/.blitz/tarballs"),
		installPath:  envOr(envInstallPath, "~/.blitz/generators"),
	}
}
