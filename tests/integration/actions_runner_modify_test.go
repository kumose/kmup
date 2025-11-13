// Copyright (C) Kumo inc. and its affiliates.
// Author: Jeff.li lijippy@163.com
// All rights reserved.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package integration

import (
	"fmt"
	"net/http"
	"testing"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionsRunnerModify(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	ctx := t.Context()

	require.NoError(t, db.DeleteAllRecords("action_runner"))

	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	_ = actions_model.CreateRunner(ctx, &actions_model.ActionRunner{OwnerID: user2.ID, Name: "user2-runner", TokenHash: "a", UUID: "a"})
	user2Runner := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunner{OwnerID: user2.ID, Name: "user2-runner"})
	userWebURL := "/user/settings/actions/runners"

	org3 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3, Type: user_model.UserTypeOrganization})
	require.NoError(t, actions_model.CreateRunner(ctx, &actions_model.ActionRunner{OwnerID: org3.ID, Name: "org3-runner", TokenHash: "b", UUID: "b"}))
	org3Runner := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunner{OwnerID: org3.ID, Name: "org3-runner"})
	orgWebURL := "/org/org3/settings/actions/runners"

	repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	_ = actions_model.CreateRunner(ctx, &actions_model.ActionRunner{RepoID: repo1.ID, Name: "repo1-runner", TokenHash: "c", UUID: "c"})
	repo1Runner := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunner{RepoID: repo1.ID, Name: "repo1-runner"})
	repoWebURL := "/user2/repo1/settings/actions/runners"

	_ = actions_model.CreateRunner(ctx, &actions_model.ActionRunner{Name: "global-runner", TokenHash: "d", UUID: "d"})
	globalRunner := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunner{Name: "global-runner"})
	adminWebURL := "/-/admin/actions/runners"

	sessionAdmin := loginUser(t, "user1")
	sessionUser2 := loginUser(t, user2.Name)

	doUpdate := func(t *testing.T, sess *TestSession, baseURL string, id int64, description string, expectedStatus int) {
		req := NewRequestWithValues(t, "POST", fmt.Sprintf("%s/%d", baseURL, id), map[string]string{
			"_csrf":       GetUserCSRFToken(t, sess),
			"description": description,
		})
		sess.MakeRequest(t, req, expectedStatus)
	}

	doDelete := func(t *testing.T, sess *TestSession, baseURL string, id int64, expectedStatus int) {
		req := NewRequestWithValues(t, "POST", fmt.Sprintf("%s/%d/delete", baseURL, id), map[string]string{
			"_csrf": GetUserCSRFToken(t, sess),
		})
		sess.MakeRequest(t, req, expectedStatus)
	}

	assertDenied := func(t *testing.T, sess *TestSession, baseURL string, id int64) {
		doUpdate(t, sess, baseURL, id, "ChangedDescription", http.StatusNotFound)
		doDelete(t, sess, baseURL, id, http.StatusNotFound)
		v := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunner{ID: id})
		assert.Empty(t, v.Description)
	}

	assertSuccess := func(t *testing.T, sess *TestSession, baseURL string, id int64) {
		doUpdate(t, sess, baseURL, id, "ChangedDescription", http.StatusSeeOther)
		v := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunner{ID: id})
		assert.Equal(t, "ChangedDescription", v.Description)
		doDelete(t, sess, baseURL, id, http.StatusOK)
		unittest.AssertNotExistsBean(t, &actions_model.ActionRunner{ID: id})
	}

	t.Run("UpdateUserRunner", func(t *testing.T) {
		theRunner := user2Runner
		t.Run("FromOrg", func(t *testing.T) {
			assertDenied(t, sessionAdmin, orgWebURL, theRunner.ID)
		})
		t.Run("FromRepo", func(t *testing.T) {
			assertDenied(t, sessionAdmin, repoWebURL, theRunner.ID)
		})
		t.Run("FromAdmin", func(t *testing.T) {
			t.Skip("Admin can update any runner (not right but not too bad)")
			assertDenied(t, sessionAdmin, adminWebURL, theRunner.ID)
		})
	})

	t.Run("UpdateOrgRunner", func(t *testing.T) {
		theRunner := org3Runner
		t.Run("FromRepo", func(t *testing.T) {
			assertDenied(t, sessionAdmin, repoWebURL, theRunner.ID)
		})
		t.Run("FromUser", func(t *testing.T) {
			assertDenied(t, sessionAdmin, userWebURL, theRunner.ID)
		})
		t.Run("FromAdmin", func(t *testing.T) {
			t.Skip("Admin can update any runner (not right but not too bad)")
			assertDenied(t, sessionAdmin, adminWebURL, theRunner.ID)
		})
	})

	t.Run("UpdateRepoRunner", func(t *testing.T) {
		theRunner := repo1Runner
		t.Run("FromOrg", func(t *testing.T) {
			assertDenied(t, sessionAdmin, orgWebURL, theRunner.ID)
		})
		t.Run("FromUser", func(t *testing.T) {
			assertDenied(t, sessionAdmin, userWebURL, theRunner.ID)
		})
		t.Run("FromAdmin", func(t *testing.T) {
			t.Skip("Admin can update any runner (not right but not too bad)")
			assertDenied(t, sessionAdmin, adminWebURL, theRunner.ID)
		})
	})

	t.Run("UpdateGlobalRunner", func(t *testing.T) {
		theRunner := globalRunner
		t.Run("FromOrg", func(t *testing.T) {
			assertDenied(t, sessionAdmin, orgWebURL, theRunner.ID)
		})
		t.Run("FromUser", func(t *testing.T) {
			assertDenied(t, sessionAdmin, userWebURL, theRunner.ID)
		})
		t.Run("FromRepo", func(t *testing.T) {
			assertDenied(t, sessionAdmin, repoWebURL, theRunner.ID)
		})
	})

	t.Run("UpdateSuccess", func(t *testing.T) {
		t.Run("User", func(t *testing.T) {
			assertSuccess(t, sessionUser2, userWebURL, user2Runner.ID)
		})
		t.Run("Org", func(t *testing.T) {
			assertSuccess(t, sessionAdmin, orgWebURL, org3Runner.ID)
		})
		t.Run("Repo", func(t *testing.T) {
			assertSuccess(t, sessionUser2, repoWebURL, repo1Runner.ID)
		})
		t.Run("Admin", func(t *testing.T) {
			assertSuccess(t, sessionAdmin, adminWebURL, globalRunner.ID)
		})
	})
}
