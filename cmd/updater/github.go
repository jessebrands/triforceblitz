package main

import (
	"context"
	"github.com/google/go-github/v68/github"
	"github.com/jessebrands/triforceblitz/internal/generator"
)

type GithubSource struct {
	client *github.Client
	Owner  string
	Repo   string
}

func NewGitHubSource(client *github.Client, owner, repo string) *GithubSource {
	return &GithubSource{
		client: client,
		Owner:  owner,
		Repo:   repo,
	}
}

func (s *GithubSource) ListAvailable(ctx context.Context) ([]Package, error) {
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
			packages = append(packages, Package{
				Version:     version,
				PublishedAt: r.GetPublishedAt().Time,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return packages, nil
}

func (s *GithubSource) Type() string {
	return "github"
}

func (s *GithubSource) Name() string {
	return s.Owner + "/" + s.Repo
}

func (s *GithubSource) String() string {
	return SourceIdentifier(s)
}
