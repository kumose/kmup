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
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIRepoBranchesPlain(t *testing.T) {
	onKmupRun(t, func(*testing.T, *url.URL) {
		repo3 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
		user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
		session := loginUser(t, user1.LowerName)

		// public only token should be forbidden
		publicOnlyToken := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopePublicOnly, auth_model.AccessTokenScopeWriteRepository)
		link, _ := url.Parse(fmt.Sprintf("/api/v1/repos/org3/%s/branches", repo3.Name)) // a plain repo
		MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(publicOnlyToken), http.StatusForbidden)

		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
		resp := MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
		bs, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var branches []*api.Branch
		assert.NoError(t, json.Unmarshal(bs, &branches))
		assert.Len(t, branches, 2)
		assert.Equal(t, "test_branch", branches[0].Name)
		assert.Equal(t, "master", branches[1].Name)

		link2, _ := url.Parse(fmt.Sprintf("/api/v1/repos/org3/%s/branches/test_branch", repo3.Name))
		MakeRequest(t, NewRequest(t, "GET", link2.String()).AddTokenAuth(publicOnlyToken), http.StatusForbidden)

		resp = MakeRequest(t, NewRequest(t, "GET", link2.String()).AddTokenAuth(token), http.StatusOK)
		bs, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)
		var branch api.Branch
		assert.NoError(t, json.Unmarshal(bs, &branch))
		assert.Equal(t, "test_branch", branch.Name)

		MakeRequest(t, NewRequest(t, "POST", link.String()).AddTokenAuth(publicOnlyToken), http.StatusForbidden)

		req := NewRequest(t, "POST", link.String()).AddTokenAuth(token)
		req.Header.Add("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(`{"new_branch_name":"test_branch2", "old_branch_name": "test_branch", "old_ref_name":"refs/heads/test_branch"}`))
		resp = MakeRequest(t, req, http.StatusCreated)
		bs, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)
		var branch2 api.Branch
		assert.NoError(t, json.Unmarshal(bs, &branch2))
		assert.Equal(t, "test_branch2", branch2.Name)
		assert.Equal(t, branch.Commit.ID, branch2.Commit.ID)

		resp = MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
		bs, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)

		branches = []*api.Branch{}
		assert.NoError(t, json.Unmarshal(bs, &branches))
		assert.Len(t, branches, 3)
		assert.Equal(t, "test_branch", branches[0].Name)
		assert.Equal(t, "test_branch2", branches[1].Name)
		assert.Equal(t, "master", branches[2].Name)

		link3, _ := url.Parse(fmt.Sprintf("/api/v1/repos/org3/%s/branches/test_branch2", repo3.Name))
		MakeRequest(t, NewRequest(t, "DELETE", link3.String()), http.StatusNotFound)
		MakeRequest(t, NewRequest(t, "DELETE", link3.String()).AddTokenAuth(publicOnlyToken), http.StatusForbidden)

		MakeRequest(t, NewRequest(t, "DELETE", link3.String()).AddTokenAuth(token), http.StatusNoContent)
		assert.NoError(t, err)
	})
}

func TestAPIRepoBranchesMirror(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo5 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 5})
	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	session := loginUser(t, user1.LowerName)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

	link, _ := url.Parse(fmt.Sprintf("/api/v1/repos/org3/%s/branches", repo5.Name)) // a mirror repo
	resp := MakeRequest(t, NewRequest(t, "GET", link.String()).AddTokenAuth(token), http.StatusOK)
	bs, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var branches []*api.Branch
	assert.NoError(t, json.Unmarshal(bs, &branches))
	assert.Len(t, branches, 2)
	assert.Equal(t, "test_branch", branches[0].Name)
	assert.Equal(t, "master", branches[1].Name)

	link2, _ := url.Parse(fmt.Sprintf("/api/v1/repos/org3/%s/branches/test_branch", repo5.Name))
	resp = MakeRequest(t, NewRequest(t, "GET", link2.String()).AddTokenAuth(token), http.StatusOK)
	bs, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var branch api.Branch
	assert.NoError(t, json.Unmarshal(bs, &branch))
	assert.Equal(t, "test_branch", branch.Name)

	req := NewRequest(t, "POST", link.String()).AddTokenAuth(token)
	req.Header.Add("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(`{"new_branch_name":"test_branch2", "old_branch_name": "test_branch", "old_ref_name":"refs/heads/test_branch"}`))
	resp = MakeRequest(t, req, http.StatusForbidden)
	bs, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, "{\"message\":\"Git Repository is a mirror.\",\"url\":\""+setting.AppURL+"api/swagger\"}", string(bs))

	resp = MakeRequest(t, NewRequest(t, "DELETE", link2.String()).AddTokenAuth(token), http.StatusForbidden)
	bs, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, "{\"message\":\"Git Repository is a mirror.\",\"url\":\""+setting.AppURL+"api/swagger\"}", string(bs))
}
