package reconciler

import (
	"reflect"
	"testing"

	"github.com/cpanato/github-sync/pkg/config"
)

func TestReconcileCollaborators(t *testing.T) {
	tests := []struct {
		name               string
		existingRepository []config.Repository
		configRepository   []config.Repository
		expectedActions    []Action
	}{
		{
			name: "add new collaborator",
			existingRepository: []config.Repository{
				{
					Name: "honk_repo",
					Collaborators: []config.Collaborator{
						{
							Username:   "honk",
							Email:      "test@example.com",
							Permission: "push",
						},
					},
				},
			},
			configRepository: []config.Repository{
				{
					Name: "honk_repo",
					Collaborators: []config.Collaborator{
						{
							Username:   "honk",
							Email:      "test@example.com",
							Permission: "push",
						},
						{
							Username:   "honk_2",
							Email:      "test_2@example.com",
							Permission: "push",
						},
					},
				},
			},
			expectedActions: []Action{&inviteCollaboratorAction{config.Collaborator{Username: "honk_2", Email: "test_2@example.com", Permission: "push"}, "honk_repo", "honk_org"}},
		},
		{
			name: "remove collaborator",
			existingRepository: []config.Repository{
				{
					Name: "honk_repo",
					Collaborators: []config.Collaborator{
						{
							Username:   "honk",
							Email:      "test@example.com",
							Permission: "push",
						},
						{
							Username:   "honk_2",
							Email:      "test_2@example.com",
							Permission: "push",
						},
					},
				},
			},
			configRepository: []config.Repository{
				{
					Name: "honk_repo",
					Collaborators: []config.Collaborator{
						{
							Username:   "honk",
							Email:      "test@example.com",
							Permission: "push",
						},
					},
				},
			},
			expectedActions: []Action{&removeCollaboratorAction{"honk_2", "honk_org", "honk_repo"}},
		},
		{
			name: "no changes",
			existingRepository: []config.Repository{
				{
					Name: "honk_repo",
					Collaborators: []config.Collaborator{
						{
							Username:   "honk",
							Email:      "test@example.com",
							Permission: "push",
						},
						{
							Username:   "honk_2",
							Email:      "test_2@example.com",
							Permission: "push",
						},
					},
				},
			},
			configRepository: []config.Repository{
				{
					Name: "honk_repo",
					Collaborators: []config.Collaborator{
						{
							Username:   "honk",
							Email:      "test@example.com",
							Permission: "push",
						},
						{
							Username:   "honk_2",
							Email:      "test_2@example.com",
							Permission: "push",
						},
					},
				},
			},
			expectedActions: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := Reconciler{
				config:       config.Config{Repositories: tc.configRepository},
				users:        usersState{byUsername: map[string]*config.User{}},
				repositories: repoCollabState{ByRepo: map[string]*config.Repository{}},
				org:          "honk_org",
			}
			for _, c := range tc.existingRepository {
				c2 := c
				r.repositories.ByRepo[c.Name] = &c2
			}

			actions, _ := r.reconcileCollaborators(false)
			if !reflect.DeepEqual(actions, tc.expectedActions) {
				t.Errorf("Expected actions: %#v\nActual actions: %#v", tc.expectedActions, actions)
			}
		})
	}
}
