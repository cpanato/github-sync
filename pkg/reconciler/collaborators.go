package reconciler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cpanato/github-sync/pkg/config"
	"github.com/google/go-github/v35/github"
)

func (r *Reconciler) reconcileCollaborators(dryRun bool) ([]Action, []error) {
	missingCollab := map[string]config.Collaborator{}

	var actions []Action // nolint: prealloc
	var errors []error   // nolint: prealloc

	for _, re := range r.repositories.ByRepo {
		for _, co := range re.Collaborators {
			missingCollab[fmt.Sprintf("%s__%s", re.Name, co.Username)] = co
		}
	}

	for _, repo := range r.config.Repositories {
		for _, c := range repo.Collaborators {
			foundCollab := false
			for _, collabState := range r.repositories.ByRepo[repo.Name].Collaborators {
				if collabState.Username == c.Username {
					delete(missingCollab, fmt.Sprintf("%s__%s", repo.Name, collabState.Username))
					foundCollab = true
					break
				}
			}

			if !foundCollab {
				actions = append(actions, &inviteCollaboratorAction{c, repo.Name, r.org})
			}
		}
	}

	for v, o := range missingCollab {
		if dryRun {
			errors = append(errors, fmt.Errorf("collaborator %s not referenced in config will be removed", o.Username))
		}

		repo := strings.Split(v, "__")[0]
		actions = append(actions, &removeCollaboratorAction{o.Username, r.org, repo})
	}

	return actions, errors
}

type inviteCollaboratorAction struct {
	config.Collaborator
	repo string
	org  string
}

func (a *inviteCollaboratorAction) Describe() string {
	return fmt.Sprintf("Invite new collaborator: %s/%s to repository: %s with **%s** permission", a.Username, a.Email, a.repo, a.Permission)
}

func (a *inviteCollaboratorAction) Perform(reconciler *Reconciler) error {
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: a.Permission,
	}

	invitation, _, err := reconciler.github.Repositories.AddCollaborator(context.Background(), a.org, a.repo, a.Username, opts)
	if err != nil {
		log.Fatalf("failed to invite new user %s: %v\n", a.Email, err.Error())
	}

	fmt.Printf("invite %+v", invitation)

	return nil
}

type removeCollaboratorAction struct {
	username string
	org      string
	repo     string
}

func (a *removeCollaboratorAction) Describe() string {
	return fmt.Sprintf("Remove collaborator from %s/%s: %s", a.org, a.repo, a.username)
}

func (a *removeCollaboratorAction) Perform(reconciler *Reconciler) error {
	_, err := reconciler.github.Organizations.RemoveOutsideCollaborator(context.Background(), a.org, a.username)
	if err != nil {
		log.Fatalf("failed to remove %s collaborator from %s/%s: %v\n", a.username, a.org, a.repo, err.Error())
	}

	return nil
}
