package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"github.com/google/go-github/v68/github"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type GitHubSource struct {
	client   *github.Client
	cacheDir string
	index    map[generator.Version]*GitHubPackage
	Owner    string
	Repo     string
}

type GitHubPackage struct {
	Version     generator.Version
	PublishedAt time.Time
	TarballUrl  string
}

// NewGitHubSource creates a new Source that sources Packages from GitHub.
func NewGitHubSource(client *github.Client, owner, repo string, cacheDir string) *GitHubSource {
	cacheDir = filepath.Join(cacheDir, "github", owner, repo)

	return &GitHubSource{
		client:   client,
		cacheDir: cacheDir,
		index:    make(map[generator.Version]*GitHubPackage),
		Owner:    owner,
		Repo:     repo,
	}
}

func (s *GitHubSource) getPackageTarball(version generator.Version) string {
	return filepath.Join(s.cacheDir, version.String()+".tar.gz")
}

func (s *GitHubSource) Update(ctx context.Context) error {
	opt := &github.ListOptions{PerPage: 100}
	for {
		releases, resp, err := s.client.Repositories.ListReleases(ctx, s.Owner, s.Repo, opt)
		if err != nil {
			return err
		}
		for _, r := range releases {
			version, err := generator.VersionFromString(r.GetTagName())
			if err != nil {
				// Tag has an invalid Triforce Blitz version, so this is not a TFB release.
				// Just skip it, this is not an error.
				continue
			}
			if pkg, ok := s.index[version]; !ok {
				s.index[version] = &GitHubPackage{
					Version:     version,
					PublishedAt: r.GetPublishedAt().Time,
					TarballUrl:  r.GetTarballURL(),
				}
			} else {
				pkg.PublishedAt = r.GetPublishedAt().Time
				pkg.TarballUrl = r.GetTarballURL()
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil
}

func (s *GitHubSource) GetAllPackages() []Package {
	var packages []Package
	for _, pkg := range s.index {
		packages = append(packages, pkg)
	}
	return packages
}

func (s *GitHubSource) GetPackage(version generator.Version) (Package, error) {
	if pkg, ok := s.index[version]; ok {
		return pkg, nil
	}
	return nil, ErrPackageNotFound
}

func (s *GitHubSource) DownloadPackage(ctx context.Context, version generator.Version) error {
	if s.IsCached(version) {
		return nil
	}
	pkg, err := s.GetPackage(version)
	if err != nil {
		return err
	}
	// Create cache directory.
	filename := s.getPackageTarball(version)
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	return pkg.Download(ctx, filename)
}

func (s *GitHubSource) UnpackPackage(ctx context.Context, version generator.Version, destination string) error {
	if !s.IsCached(version) {
		if err := s.DownloadPackage(ctx, version); err != nil {
			return err
		}
	}
	filename := s.getPackageTarball(version)
	// We got our tarball, open a stream to it.
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	archive, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer archive.Close()
	reader := tar.NewReader(archive)
	// Ensure the destination exists and begin reading.
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(target, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			file, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, reader); err != nil {
				file.Close()
				return err
			}
			file.Close()

		default:
			continue
		}
	}
	return nil
}

func (s *GitHubSource) IsCached(version generator.Version) bool {
	if _, err := os.Stat(s.getPackageTarball(version)); err == nil {
		return true
	}
	return false
}

func (s *GitHubSource) Type() string {
	return "github"
}

func (s *GitHubSource) Name() string {
	return s.Owner + "/" + s.Repo
}

func (s *GitHubSource) String() string {
	return SourceIdentifier(s)
}

func (p *GitHubPackage) Download(ctx context.Context, destination string) error {
	f, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(p.TarballUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("file download failed: " + resp.Status)
	}
	_, err = io.Copy(f, resp.Body)
	return err
}

func (p *GitHubPackage) GetVersion() generator.Version {
	return p.Version
}

func (p *GitHubPackage) GetPublishedAt() time.Time {
	return p.PublishedAt
}
