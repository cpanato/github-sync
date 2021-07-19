package reconciler

import (
	"context"

	"github.com/google/go-github/v35/github"

	"github.com/cpanato/github-sync/pkg/config"
)

type repoCollabState struct {
	ByRepo map[string]*config.Repository
}

func (c *repoCollabState) init(gh *github.Client, org string) error {
	c.ByRepo = map[string]*config.Repository{}
	var repos []string

	page := 1
	for {
		opts := &github.RepositoryListByOrgOptions{
			Type: "all",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}

		repositories, resp, err := gh.Repositories.ListByOrg(context.Background(), org, opts)
		if err != nil {
			return err
		}

		for _, repo := range repositories {
			repos = append(repos, repo.GetName())
		}

		if resp.NextPage == 0 {
			break
		}
		page++
	}

	for _, repo := range repos {
		var assets []config.Collaborator
		page := 1
		for {
			opts := &github.ListCollaboratorsOptions{
				Affiliation: "outside",
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			}

			users, resp, err := gh.Repositories.ListCollaborators(context.Background(), org, repo, opts)
			if err != nil {
				return err
			}

			for _, user := range users {
				appendCollab := config.Collaborator{
					Username: user.GetLogin(),
					Email:    user.GetEmail(),
				}

				assets = append(assets, appendCollab)
			}
			if resp.NextPage == 0 {
				break
			}
			page++
		}

		c.ByRepo[repo] = &config.Repository{
			Name:          repo,
			Collaborators: assets,
		}
	}

	return nil
}
