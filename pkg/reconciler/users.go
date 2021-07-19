package reconciler

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v35/github"

	"github.com/cpanato/github-sync/pkg/config"
)

func (r *Reconciler) reconcileUsers(dryRun bool) ([]Action, []error) {
	missingUsers := map[string]*config.User{}

	var actions []Action // nolint: prealloc
	var errors []error   // nolint: prealloc

	for _, c := range r.users.byUsername {
		missingUsers[c.Username] = c
	}

	for _, c := range r.config.Users {
		if o, ok := r.users.byUsername[c.Username]; ok {
			delete(missingUsers, o.Username)
		} else {
			actions = append(actions, &inviteUserAction{c, r.org})
		}
	}

	for _, o := range missingUsers {
		if dryRun {
			errors = append(errors, fmt.Errorf("user %s not referenced in config will be removed", o.Username))
		}
		actions = append(actions, &removeUserAction{o.Username, r.org})
	}

	return actions, errors
}

type inviteUserAction struct {
	config.User
	org string
}

func (a *inviteUserAction) Describe() string {
	return fmt.Sprintf("Invite new user: %s/%s", a.Username, a.User)
}

func (a *inviteUserAction) Perform(reconciler *Reconciler) error {
	invite := &github.CreateOrgInvitationOptions{
		Email:  &a.Email,
		Role:   &a.Role,
		TeamID: []int64{},
	}

	invitation, _, err := reconciler.github.Organizations.CreateOrgInvitation(context.Background(), a.org, invite)
	if err != nil {
		log.Fatalf("failed to invite new user %s: %v\n", a.Email, err.Error())
	}

	newUser := &config.User{
		Email:    invitation.GetEmail(),
		Username: invitation.GetLogin(),
		Role:     invitation.GetRole(),
	}

	reconciler.users.byUsername[a.Username] = newUser
	return nil
}

type removeUserAction struct {
	username string
	org      string
}

func (a *removeUserAction) Describe() string {
	return fmt.Sprintf("Remove user from %s: %s", a.org, a.username)
}

func (a *removeUserAction) Perform(reconciler *Reconciler) error {
	_, err := reconciler.github.Organizations.RemoveMember(context.Background(), a.org, a.username)
	if err != nil {
		log.Fatalf("failed to remove member from all teams user %s: %v\n", a.username, err.Error())
	}

	_, err = reconciler.github.Organizations.RemoveOrgMembership(context.Background(), a.username, a.org)
	if err != nil {
		log.Fatalf("failed to remove member from the org %s: %v\n", a.username, err.Error())
	}

	return nil
}
