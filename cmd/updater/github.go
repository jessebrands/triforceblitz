package main

import (
	"context"
	"github.com/google/go-github/v68/github"
	"github.com/jessebrands/triforceblitz/internal/generator"
	"time"
)

type GitHubSource struct {
	client *github.Client
	Owner  string
	Repo   string
}

type GitHubPackage struct {
	Version     generator.Version
	PublishedAt time.Time
	TarballUrl  string
}

func NewGitHubSource(client *github.Client, owner, repo string) *GitHubSource {
	return &GitHubSource{
		client: client,
		Owner:  owner,
		Repo:   repo,
	}
}

func (s *GitHubSource) ListAvailable(ctx context.Context) ([]Package, error) {
	var packages []Package
	opt := &github.ListOptions{PerPage: 100}
	for {
		releases, resp, err := s.client.Repositories.ListReleases(ctx, s.Owner, s.Repo, opt)
		if err != nil {
			return packages, err
		}
		for _, r := range releases {
			version, err := generator.VersionFromString(r.GetTagName())
			if err != nil {
				// If the tag isn't a valid Generator version, that means that this is not a generator.
				// No need to report an error, just continue.
				continue
			}
			packages = append(packages, &GitHubPackage{
				Version:     version,
				PublishedAt: r.GetPublishedAt().Time,
				TarballUrl:  r.GetTarballURL(),
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return packages, nil
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
