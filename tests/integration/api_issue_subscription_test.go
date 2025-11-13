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

func TestAPIIssueSubscriptions(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	issue1 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
	issue2 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2})
	issue3 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 3})
	issue4 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 4})
	issue5 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 8})

	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: issue1.PosterID})

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	testSubscription := func(issue *issues_model.Issue, isWatching bool) {
		issueRepo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: issue.RepoID})

		req := NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/subscriptions/check", issueRepo.OwnerName, issueRepo.Name, issue.Index)).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)
		wi := new(api.WatchInfo)
		DecodeJSON(t, resp, wi)

		assert.Equal(t, isWatching, wi.Subscribed)
		assert.Equal(t, !isWatching, wi.Ignored)
		assert.Equal(t, issue.APIURL(t.Context())+"/subscriptions", wi.URL)
		assert.EqualValues(t, issue.CreatedUnix, wi.CreatedAt.Unix())
		assert.Equal(t, issueRepo.APIURL(), wi.RepositoryURL)
	}

	testSubscription(issue1, true)
	testSubscription(issue2, true)
	testSubscription(issue3, true)
	testSubscription(issue4, false)
	testSubscription(issue5, false)

	issue1Repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: issue1.RepoID})
	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/subscriptions/%s", issue1Repo.OwnerName, issue1Repo.Name, issue1.Index, owner.Name)
	req := NewRequest(t, "DELETE", urlStr).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusCreated)
	testSubscription(issue1, false)

	req = NewRequest(t, "DELETE", urlStr).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusOK)
	testSubscription(issue1, false)

	issue5Repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: issue5.RepoID})
	urlStr = fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/subscriptions/%s", issue5Repo.OwnerName, issue5Repo.Name, issue5.Index, owner.Name)
	req = NewRequest(t, "PUT", urlStr).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusCreated)
	testSubscription(issue5, true)

	req = NewRequest(t, "PUT", urlStr).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusOK)
	testSubscription(issue5, true)
}
