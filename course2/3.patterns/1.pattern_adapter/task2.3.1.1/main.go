package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "your_access_token"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	var g Githuber
	g = NewGithubAdapter(client)
	fmt.Println(g.GetGists(context.Background(), "ptflp"))
	fmt.Println(g.GetRepos(context.Background(), "ptflp"))
}

type RepoLister interface {
	List(ctx context.Context, username string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
}

type GistLister interface {
	List(ctx context.Context, username string, opt *github.GistListOptions) ([]*github.Gist, *github.Response, error)
}
type Githuber interface {
	GetGists(ctx context.Context, username string) ([]Item, error)
	GetRepos(ctx context.Context, username string) ([]Item, error) // opt :=&github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 1000}}
}
type GithubAdapter struct {
	RepoList RepoLister
	GistList GistLister
}

func (a *GithubAdapter) GetGists(ctx context.Context, username string) ([]Item, error) {
	g, _, err := a.GistList.List(ctx, username, nil)
	if err != nil {
		return nil, err
	}
	res := make([]Item, len(g))
	for i := range g {
		res[i] = Item{
			Title:       g[i].GetID(),
			Description: g[i].GetDescription(),
			Link:        g[i].GetHTMLURL(),
		}
	}
	return res, nil
}

func (a *GithubAdapter) GetRepos(ctx context.Context, username string) ([]Item, error) {
	r, _, err := a.RepoList.List(ctx, username, nil)
	if err != nil {
		return nil, err
	}
	res := make([]Item, len(r))
	for i := range r {
		res[i] = Item{
			Title:       r[i].GetName(),
			Description: r[i].GetDescription(),
			Link:        r[i].GetHTMLURL(),
		}
	}
	return res, nil
}

func NewGithubAdapter(githubClient *github.Client) *GithubAdapter {
	g := &GithubAdapter{
		RepoList: githubClient.Repositories,
		GistList: githubClient.Gists,
	}
	return g
}

type Item struct {
	Title       string
	Description string
	Link        string
}
