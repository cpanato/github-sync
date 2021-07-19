package github

import (
	"context"
	"errors"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v35/github"
)

func NewClient(githubAccessToken string) (*github.Client, error) {
	if githubAccessToken == "" {
		return nil, errors.New("missing GitHub Access Token")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}
