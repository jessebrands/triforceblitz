package main

import (
	"context"
	"errors"
	"github.com/google/go-github/v68/github"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"time"
)

type GitHubSource struct {
	client   *github.Client
	index    map[generator.Version]*GitHubPackage
	Owner    string
	Repo     string
}

type GitHubPackage struct {
	Version     generator.Version
	PublishedAt time.Time
	TarballUrl  string
}

func NewGitHubSource(client *github.Client, owner, repo string) *GitHubSource {
	return &GitHubSource{
		client:   client,
		index:    make(map[generator.Version]*GitHubPackage),
		Owner:    owner,
		Repo:     repo,
	}
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
	return nil, errors.New("package not found")
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

func (p *GitHubPackage) GetVersion() generator.Version {
	return p.Version
}

func (p *GitHubPackage) GetPublishedAt() time.Time {
	return p.PublishedAt
}
