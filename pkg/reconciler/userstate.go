package reconciler

import (
	"context"

	"github.com/google/go-github/v35/github"

	"github.com/cpanato/github-sync/pkg/config"
)

type usersState struct {
	byUsername map[string]*config.User
}

func (c *usersState) init(gh *github.Client, org string) error {
	c.byUsername = map[string]*config.User{}

	var assets []config.User
	page := 1
	for {
		opts := &github.ListMembersOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}
		users, resp, err := gh.Organizations.ListMembers(context.Background(), org, opts)
		if err != nil {
			return err
		}
		for _, user := range users {
			userInfo, _, err := gh.Users.GetByID(context.Background(), user.GetID())
			if err != nil {
				return err
			}
			appendUser := config.User{
				Username: user.GetLogin(),
				Email:    userInfo.GetEmail(),
			}

			assets = append(assets, appendUser)
		}
		if resp.NextPage == 0 {
			break
		}
		page++
	}

	for _, user := range assets {
		user2 := user
		c.byUsername[user.Username] = &user2
	}

	return nil
}
