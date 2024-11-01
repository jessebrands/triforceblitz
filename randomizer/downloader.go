package randomizer

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v66/github"
)

type Service struct {
	// GitHub client
	client *github.Client
	repository repository

	// Path to download release tarballs to.
	downloadPath string 

	// Path to install releases to.
	installPath string 
}

func NewService(config *Config) *Service {
	return &Service{
		client: github.NewClient(nil),
		repository: config.repository,
		downloadPath: config.downloadPath,
		installPath: config.installPath,
	}
}

type Release struct {
	Version  string
	AssetURL string
}

type Downloader struct {
	client  *github.Client
	owner   string
	repo    string
	destDir string
}

// NewDownloader creates a downloader that can query the given repository as specified by
// owner and repo. The queried releases can then be downloaded to the directory specified by
// dest.
func NewDownloader(owner string, repo string, dest string) *Downloader {
	return &Downloader{
		owner:   owner,
		repo:    repo,
		destDir: dest,
		client:  github.NewClient(nil),
	}
}

func (d *Downloader) GetAvailableReleases() ([]*Release, error) {
	results := []*Release{}
	opts := &github.ListOptions{
		PerPage: 100,
	}
	releases, resp, err := d.client.Repositories.ListReleases(context.Background(), d.owner, d.repo, opts)
	if err != nil {
		return results, err
	}
	if resp.StatusCode != 200 {
		return results, fmt.Errorf("Failed to get releases, HTTP Error %d", resp.StatusCode)
	}
	for _, rel := range releases {
		results = append(results, &Release{
			Version:  strings.TrimSpace(rel.GetName()),
			AssetURL: rel.GetTarballURL(),
		})
	}
	return results, err
}

func (d *Downloader) Download(release *Release) (string, error) {
	outPath := fmt.Sprintf("%s/tarballs/%s.tar.gz", d.destDir, release.Version)
	file, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(release.AssetURL)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(file, resp.Body)
	return outPath, err
}

func (d *Downloader) Install(tarball string) error {
	file, err := os.Open(tarball)
	defer file.Close()
	if err != nil {
		return err
	}
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()
	r := tar.NewReader(gz)
	// Grab the first file and check the header. If it's just a directory then
	// it is probably the parent directory. We want to skip this file so save 
	// the name and trim it as a prefix of subsequent headers. Too bad there
	// isn't a better way of doing this.
	header, err := r.Next() 
	if err != nil {
		return err 
	}
	if header == nil {
		return fmt.Errorf("Got unexpected empty tarball header")
	}
	dest := "/home/bee/releases/out"
	parent := ""
	os.MkdirAll(dest, 0744)
	log.Println("Extracting files to", dest)
	for {
		header, err := r.Next()
		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		filename := strings.TrimPrefix(header.Name, parent)
		target := filepath.Join(dest, filename)

		switch header.Typeflag { 
		case tar.TypeDir: 
			if len(parent) == 0 {
				parent = header.Name 
				log.Println("Ignoring prefix", parent)
				continue
			}
			log.Println("Creating directory:", target)
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err 
				}
			}

		case tar.TypeReg: 
			log.Println("Extracting file:", target)
			f, err := os.OpenFile(target, os.O_CREATE | os.O_RDWR, os.FileMode(header.Mode))
			defer f.Close()
			if err != nil {
				return err 
			}
			if _, err := io.Copy(f, r); err != nil {
				return err 
			}
		}
	}
}

