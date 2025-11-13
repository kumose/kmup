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

func TestAPIPinIssue(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	assert.NoError(t, unittest.LoadFixtures())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{RepoID: repo.ID})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	// Pin the Issue
	req := NewRequest(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", repo.OwnerName, repo.Name, issue.Index)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Check if the Issue is pinned
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d", repo.OwnerName, repo.Name, issue.Index))
	resp := MakeRequest(t, req, http.StatusOK)
	var issueAPI api.Issue
	DecodeJSON(t, resp, &issueAPI)
	assert.Equal(t, 1, issueAPI.PinOrder)
}

func TestAPIUnpinIssue(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	assert.NoError(t, unittest.LoadFixtures())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{RepoID: repo.ID})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	// Pin the Issue
	req := NewRequest(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", repo.OwnerName, repo.Name, issue.Index)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Check if the Issue is pinned
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d", repo.OwnerName, repo.Name, issue.Index))
	resp := MakeRequest(t, req, http.StatusOK)
	var issueAPI api.Issue
	DecodeJSON(t, resp, &issueAPI)
	assert.Equal(t, 1, issueAPI.PinOrder)

	// Unpin the Issue
	req = NewRequest(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", repo.OwnerName, repo.Name, issue.Index)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Check if the Issue is no longer pinned
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d", repo.OwnerName, repo.Name, issue.Index))
	resp = MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, &issueAPI)
	assert.Equal(t, 0, issueAPI.PinOrder)
}

func TestAPIMoveIssuePin(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	assert.NoError(t, unittest.LoadFixtures())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{RepoID: repo.ID})
	issue2 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2, RepoID: repo.ID})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	// Pin the first Issue
	req := NewRequest(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", repo.OwnerName, repo.Name, issue.Index)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Check if the first Issue is pinned at position 1
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d", repo.OwnerName, repo.Name, issue.Index))
	resp := MakeRequest(t, req, http.StatusOK)
	var issueAPI api.Issue
	DecodeJSON(t, resp, &issueAPI)
	assert.Equal(t, 1, issueAPI.PinOrder)

	// Pin the second Issue
	req = NewRequest(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", repo.OwnerName, repo.Name, issue2.Index)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Move the first Issue to position 2
	req = NewRequest(t, "PATCH", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin/2", repo.OwnerName, repo.Name, issue.Index)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Check if the first Issue is pinned at position 2
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d", repo.OwnerName, repo.Name, issue.Index))
	resp = MakeRequest(t, req, http.StatusOK)
	var issueAPI3 api.Issue
	DecodeJSON(t, resp, &issueAPI3)
	assert.Equal(t, 2, issueAPI3.PinOrder)

	// Check if the second Issue is pinned at position 1
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d", repo.OwnerName, repo.Name, issue2.Index))
	resp = MakeRequest(t, req, http.StatusOK)
	var issueAPI4 api.Issue
	DecodeJSON(t, resp, &issueAPI4)
	assert.Equal(t, 1, issueAPI4.PinOrder)
}

func TestAPIListPinnedIssues(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	assert.NoError(t, unittest.LoadFixtures())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{RepoID: repo.ID})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)

	// Pin the Issue
	req := NewRequest(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", repo.OwnerName, repo.Name, issue.Index)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Check if the Issue is in the List
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/issues/pinned", repo.OwnerName, repo.Name))
	resp := MakeRequest(t, req, http.StatusOK)
	var issueList []api.Issue
	DecodeJSON(t, resp, &issueList)

	assert.Len(t, issueList, 1)
	assert.Equal(t, issue.ID, issueList[0].ID)
}

func TestAPIListPinnedPullrequests(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	assert.NoError(t, unittest.LoadFixtures())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	req := NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/pulls/pinned", repo.OwnerName, repo.Name))
	resp := MakeRequest(t, req, http.StatusOK)
	var prList []api.PullRequest
	DecodeJSON(t, resp, &prList)

	assert.Empty(t, prList)
}

func TestAPINewPinAllowed(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	req := NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/new_pin_allowed", owner.Name, repo.Name))
	resp := MakeRequest(t, req, http.StatusOK)

	var newPinsAllowed api.NewIssuePinsAllowed
	DecodeJSON(t, resp, &newPinsAllowed)

	assert.True(t, newPinsAllowed.Issues)
	assert.True(t, newPinsAllowed.PullRequests)
}
