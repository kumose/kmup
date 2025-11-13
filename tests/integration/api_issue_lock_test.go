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

	auth_model "github.com/kumose/kmup/models/auth"
	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPILockIssue(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	t.Run("Lock", func(t *testing.T) {
		issueBefore := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
		assert.False(t, issueBefore.IsLocked)
		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: issueBefore.RepoID})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
		urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/lock", owner.Name, repo.Name, issueBefore.Index)

		session := loginUser(t, owner.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

		// check lock issue
		req := NewRequestWithJSON(t, "PUT", urlStr, api.LockIssueOption{Reason: "Spam"}).AddTokenAuth(token)
		MakeRequest(t, req, http.StatusNoContent)
		issueAfter := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
		assert.True(t, issueAfter.IsLocked)

		// check with other user
		user34 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 34})
		session34 := loginUser(t, user34.Name)
		token34 := getTokenForLoggedInUser(t, session34, auth_model.AccessTokenScopeAll)
		req = NewRequestWithJSON(t, "PUT", urlStr, api.LockIssueOption{Reason: "Spam"}).AddTokenAuth(token34)
		MakeRequest(t, req, http.StatusForbidden)
	})

	t.Run("Unlock", func(t *testing.T) {
		issueBefore := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: issueBefore.RepoID})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
		urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/lock", owner.Name, repo.Name, issueBefore.Index)

		session := loginUser(t, owner.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

		lockReq := NewRequestWithJSON(t, "PUT", urlStr, api.LockIssueOption{Reason: "Spam"}).AddTokenAuth(token)
		MakeRequest(t, lockReq, http.StatusNoContent)

		// check unlock issue
		req := NewRequest(t, "DELETE", urlStr).AddTokenAuth(token)
		MakeRequest(t, req, http.StatusNoContent)
		issueAfter := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
		assert.False(t, issueAfter.IsLocked)

		// check with other user
		user34 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 34})
		session34 := loginUser(t, user34.Name)
		token34 := getTokenForLoggedInUser(t, session34, auth_model.AccessTokenScopeAll)
		req = NewRequest(t, "DELETE", urlStr).AddTokenAuth(token34)
		MakeRequest(t, req, http.StatusForbidden)
	})
}
