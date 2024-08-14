package main

import (
	"context"
	"testing"

	"github.com/google/go-github/v53/github"
)

type MockGistLister struct {
	gists []*github.Gist
}

func (m *MockGistLister) List(ctx context.Context, username string, opt *github.GistListOptions) ([]*github.Gist, *github.Response, error) {

	return m.gists, nil, nil
}

type MockRepoLister struct {
	reps []*github.Repository
}

func (m *MockRepoLister) List(ctx context.Context, username string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return m.reps, nil, nil
}

func TestRepoList(t *testing.T) {
	mockRepoService := &MockRepoLister{
		reps: []*github.Repository{
			{
				Name:        github.String("1"),
				Description: github.String("Test Repo 1"),
				CloneURL:    github.String("http:"),
			},
			{
				Name:        github.String("2"),
				Description: github.String("Test Repo 2"),
				CloneURL:    github.String("http1:"),
			},
			{
				Name:        github.String("3"),
				Description: github.String("Test Repo 3"),
				CloneURL:    github.String("http2:"),
			},
			{
				Name:        github.String("4"),
				Description: github.String("Test Repo 4"),
				CloneURL:    github.String("http3:"),
			},
		},
	}

	adapter := GithubAdapter{RepoList: mockRepoService}
	ctx := context.Background()
	repos, err := adapter.GetRepos(ctx, "test")

	if err != nil {
		t.Errorf("expected no error, got error: %v", err)
	}

	if len(repos) != 4 {
		t.Fatalf("expected 2 gists, got %d", len(repos))
	}

	if repos[0].Title != "1" || repos[1].Description != "Test Repo 2" {
		t.Fatalf("unexpected gists: %+v", repos)
	}
}

func TestGistList(t *testing.T) {
	mockGistService := &MockGistLister{
		gists: []*github.Gist{
			{
				ID:          github.String("1"),
				Description: github.String("Test Gist 1"),
				GitPullURL:  github.String("http:"),
			},
			{
				ID:          github.String("2"),
				Description: github.String("Test Gist 2"),
				GitPullURL:  github.String("http1:"),
			},
			{
				ID:          github.String("3"),
				Description: github.String("Test Gist 3"),
				GitPullURL:  github.String("http2:"),
			},
			{
				ID:          github.String("4"),
				Description: github.String("Test Gist 4"),
				GitPullURL:  github.String("http3g:"),
			},
		},
	}

	adapter := GithubAdapter{GistList: mockGistService}
	ctx := context.Background()
	gists, err := adapter.GetGists(ctx, "test")

	if err != nil {
		t.Errorf("expected no error, got error: %v", err)
	}

	if len(gists) != 4 {
		t.Fatalf("expected 2 gists, got %d", len(gists))
	}

	if gists[0].Description != "Test Gist 1" || gists[1].Description != "Test Gist 2" {
		t.Fatalf("unexpected gists: %+v", gists)
	}
}
