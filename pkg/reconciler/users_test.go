package reconciler

import (
	"reflect"
	"testing"

	"github.com/cpanato/github-sync/pkg/config"
)

func TestReconcileUsers(t *testing.T) {
	tests := []struct {
		name             string
		existingUsers    []config.User
		configUsers      []config.User
		expectedActions  []Action
		expectedErrCount int
	}{
		{
			name:            "add new user",
			existingUsers:   []config.User{{Username: "honk", Email: "test@example.com", Role: "direct_member"}},
			configUsers:     []config.User{{Username: "honk", Email: "test@example.com", Role: "direct_member"}, {Username: "honk_2", Email: "test_2@example.com", Role: "admin"}},
			expectedActions: []Action{&inviteUserAction{config.User{Username: "honk_2", Email: "test_2@example.com", Role: "admin"}, "honk_org"}},
		},
		{
			name:            "remove user",
			existingUsers:   []config.User{{Username: "honk", Email: "test@example.com", Role: "direct_member"}, {Username: "honk_2", Email: "test_2@example.com", Role: "admin"}},
			configUsers:     []config.User{{Username: "honk", Email: "test@example.com", Role: "direct_member"}},
			expectedActions: []Action{&removeUserAction{"honk_2", "honk_org"}},
		},
		{
			name:            "no changes",
			existingUsers:   []config.User{{Username: "honk", Email: "test@example.com", Role: "direct_member"}, {Username: "honk_2", Email: "test_2@example.com", Role: "admin"}},
			configUsers:     []config.User{{Username: "honk", Email: "test@example.com", Role: "direct_member"}, {Username: "honk_2", Email: "test_2@example.com", Role: "admin"}},
			expectedActions: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := Reconciler{
				config:       config.Config{Users: tc.configUsers},
				users:        usersState{byUsername: map[string]*config.User{}},
				repositories: repoCollabState{ByRepo: map[string]*config.Repository{}},
				org:          "honk_org",
			}
			for _, c := range tc.existingUsers {
				c2 := c
				r.users.byUsername[c.Username] = &c2
			}

			actions, errs := r.reconcileUsers(false)
			if !reflect.DeepEqual(actions, tc.expectedActions) {
				t.Errorf("Expected actions: %#v\nActual actions: %#v", tc.expectedActions, actions)
			}
			if len(errs) != tc.expectedErrCount {
				t.Errorf("Expected %d errors, but got %d: %v", tc.expectedErrCount, len(errs), errs)
			}
		})
	}
}
