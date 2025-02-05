package main

import (
	"context"
	"github.com/google/go-github/v68/github"
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
			packages = append(packages, Package{
				Version:     r.GetTagName(),
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
