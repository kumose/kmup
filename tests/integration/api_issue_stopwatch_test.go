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

func TestAPIListStopWatches(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository, auth_model.AccessTokenScopeReadUser)
	req := NewRequest(t, "GET", "/api/v1/user/stopwatches").
		AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusOK)
	var apiWatches []*api.StopWatch
	DecodeJSON(t, resp, &apiWatches)
	stopwatch := unittest.AssertExistsAndLoadBean(t, &issues_model.Stopwatch{UserID: owner.ID})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: stopwatch.IssueID})
	if assert.Len(t, apiWatches, 1) {
		assert.Equal(t, stopwatch.CreatedUnix.AsTime().Unix(), apiWatches[0].Created.Unix())
		assert.Equal(t, issue.Index, apiWatches[0].IssueIndex)
		assert.Equal(t, issue.Title, apiWatches[0].IssueTitle)
		assert.Equal(t, repo.Name, apiWatches[0].RepoName)
		assert.Equal(t, repo.OwnerName, apiWatches[0].RepoOwnerName)
		assert.Positive(t, apiWatches[0].Seconds)
	}
}

func TestAPIStopStopWatches(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2})
	_ = issue.LoadRepo(t.Context())
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: issue.Repo.OwnerID})
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	req := NewRequestf(t, "POST", "/api/v1/repos/%s/%s/issues/%d/stopwatch/stop", owner.Name, issue.Repo.Name, issue.Index).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusCreated)
	MakeRequest(t, req, http.StatusConflict)
}

func TestAPICancelStopWatches(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
	_ = issue.LoadRepo(t.Context())
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: issue.Repo.OwnerID})
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	req := NewRequestf(t, "DELETE", "/api/v1/repos/%s/%s/issues/%d/stopwatch/delete", owner.Name, issue.Repo.Name, issue.Index).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)
	MakeRequest(t, req, http.StatusConflict)
}

func TestAPIStartStopWatches(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 3})
	_ = issue.LoadRepo(t.Context())
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: issue.Repo.OwnerID})
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	req := NewRequestf(t, "POST", "/api/v1/repos/%s/%s/issues/%d/stopwatch/start", owner.Name, issue.Repo.Name, issue.Index).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusCreated)
	MakeRequest(t, req, http.StatusConflict)
}
