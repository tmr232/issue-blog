package main

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/gregjones/httpcache"
	"github.com/tmr232/goat"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/tmr232/goat/cmd/goater

func getIssues(ctx context.Context, client *github.Client, owner, repo string) ([]*github.Issue, error) {
	// TODO: add pagination support
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, nil)
	return issues, err
}
func newClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := &http.Client{
		Transport: &oauth2.Transport{
			Base:   httpcache.NewMemoryCacheTransport(),
			Source: ts,
		},
	}

	client := github.NewClient(tc)
	return client
}

func renderPost(issue *github.Issue) (string, error) {
	if issue.Title == nil {
		return "", errors.New("Posts need a title!")
	}

	if issue.Body == nil {
		return "", errors.New("Posts need a body!")
	}

	title := *issue.Title
	body := *issue.Body
	date := issue.UpdatedAt.Format("2006-01-02")

	return "---\n" + "title: " + title + "\ndate: " + date + "\n---\n\n" + body, nil
}

func splitName(repoName string) (owner, repo string) {
	parts := strings.Split(repoName, "/")
	owner = parts[0]
	repo = parts[1]
	return owner, repo
}

func generateBlog(ownerAndRepo string, contentDir string, token string) error {
	goat.Flag(ownerAndRepo).Name("repo")

	client := newClient(token)

	ctx := context.Background()

	owner, repo := splitName(ownerAndRepo)

	issues, err := getIssues(ctx, client, owner, repo)
	if err != nil {
		return err
	}

	for _, issue := range issues {
		post, _ := renderPost(issue)
		err = ioutil.WriteFile(filepath.Join(contentDir, *issue.Title+".md"), []byte(post), 0606)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	goat.Run(generateBlog)
	// generateBlog("tmr232/issue-blog", filepath.Join("hugo", "content", "posts"), "token")
}
