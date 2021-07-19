package reconciler

import (
	"fmt"
	"log"

	"github.com/google/go-github/v35/github"

	"github.com/cpanato/github-sync/pkg/config"
)

type Reconciler struct {
	github       *github.Client
	config       config.Config
	users        usersState
	repositories repoCollabState
	org          string
}

func New(gh *github.Client, cfg config.Config, org string) *Reconciler {
	return &Reconciler{
		github:       gh,
		config:       cfg,
		users:        usersState{},
		repositories: repoCollabState{},
		org:          org,
	}
}

func (r *Reconciler) Reconcile(dryRun bool) error {
	if err := r.users.init(r.github, r.org); err != nil {
		return fmt.Errorf("failed to get initial user state: %v", err)
	}

	if err := r.repositories.init(r.github, r.org); err != nil {
		return fmt.Errorf("failed to get initial collaborator state: %v", err)
	}

	var actions []Action
	var errors []error
	a, e := r.reconcileUsers(dryRun)
	actions = append(actions, a...)
	errors = append(errors, e...)

	a, e = r.reconcileCollaborators(dryRun)
	actions = append(actions, a...)
	errors = append(errors, e...)

	failed := false
	if len(errors) > 0 {
		log.Printf("This configuration cannot be applied against the current reality:")
		failed = true
	}

	for i, e := range errors {
		log.Printf("Error %d: %v.\n", i+1, e)
	}

	if !dryRun && failed {
		dryRun = true
		log.Println("We will not execute anything due to errors, but this what we would've done:")
	} else if dryRun {
		log.Println("In dry run mode so taking no action, but this is what we would've done:")
	}

	if len(actions) > 0 {
		for i, a := range actions {
			log.Printf("Step %d: %s.\n", i+1, a.Describe())
			if !dryRun {
				if err := a.Perform(r); err != nil {
					log.Printf("Failed: %v.\n", err)
				}
			}
		}
	} else {
		log.Println("Nothing to do.")
	}

	if failed {
		return fmt.Errorf("there were configuration errors")
	}

	return nil
}

type Action interface {
	Describe() string
	Perform(reconciler *Reconciler) error
}
