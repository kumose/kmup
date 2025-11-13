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
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIGitHooks(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.DisableGitHooks, false)()

	const testHookContent = `#!/bin/bash
echo "TestGitHookScript"
`

	t.Run("ListGitHooks", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 37})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		// user1 is an admin user
		session := loginUser(t, "user1")
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/hooks/git", owner.Name, repo.Name).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)
		var apiGitHooks []*api.GitHook
		DecodeJSON(t, resp, &apiGitHooks)
		assert.Len(t, apiGitHooks, 3)
		for _, apiGitHook := range apiGitHooks {
			if apiGitHook.Name == "pre-receive" {
				assert.True(t, apiGitHook.IsActive)
				assert.Equal(t, testHookContent, apiGitHook.Content)
			} else {
				assert.False(t, apiGitHook.IsActive)
				assert.Empty(t, apiGitHook.Content)
			}
		}
	})

	t.Run("NoGitHooks", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		// user1 is an admin user
		session := loginUser(t, "user1")
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/hooks/git", owner.Name, repo.Name).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)
		var apiGitHooks []*api.GitHook
		DecodeJSON(t, resp, &apiGitHooks)
		assert.Len(t, apiGitHooks, 3)
		for _, apiGitHook := range apiGitHooks {
			assert.False(t, apiGitHook.IsActive)
			assert.Empty(t, apiGitHook.Content)
		}
	})

	t.Run("ListGitHooksNoAccess", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		session := loginUser(t, owner.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/hooks/git", owner.Name, repo.Name).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusForbidden)
	})

	t.Run("GetGitHook", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 37})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		// user1 is an admin user
		session := loginUser(t, "user1")
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/hooks/git/pre-receive", owner.Name, repo.Name).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)
		var apiGitHook *api.GitHook
		DecodeJSON(t, resp, &apiGitHook)
		assert.True(t, apiGitHook.IsActive)
		assert.Equal(t, testHookContent, apiGitHook.Content)
	})
	t.Run("GetGitHookNoAccess", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		session := loginUser(t, owner.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/hooks/git/pre-receive", owner.Name, repo.Name).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusForbidden)
	})

	t.Run("EditGitHook", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		// user1 is an admin user
		session := loginUser(t, "user1")
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

		urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/hooks/git/pre-receive",
			owner.Name, repo.Name)
		req := NewRequestWithJSON(t, "PATCH", urlStr, &api.EditGitHookOption{
			Content: testHookContent,
		}).AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)
		var apiGitHook *api.GitHook
		DecodeJSON(t, resp, &apiGitHook)
		assert.True(t, apiGitHook.IsActive)
		assert.Equal(t, testHookContent, apiGitHook.Content)

		req = NewRequestf(t, "GET", "/api/v1/repos/%s/%s/hooks/git/pre-receive", owner.Name, repo.Name).
			AddTokenAuth(token)
		resp = MakeRequest(t, req, http.StatusOK)
		var apiGitHook2 *api.GitHook
		DecodeJSON(t, resp, &apiGitHook2)
		assert.True(t, apiGitHook2.IsActive)
		assert.Equal(t, testHookContent, apiGitHook2.Content)
	})

	t.Run("EditGitHookNoAccess", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		session := loginUser(t, owner.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
		urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/hooks/git/pre-receive", owner.Name, repo.Name)
		req := NewRequestWithJSON(t, "PATCH", urlStr, &api.EditGitHookOption{
			Content: testHookContent,
		}).AddTokenAuth(token)
		MakeRequest(t, req, http.StatusForbidden)
	})

	t.Run("DeleteGitHook", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 37})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		// user1 is an admin user
		session := loginUser(t, "user1")
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

		req := NewRequestf(t, "DELETE", "/api/v1/repos/%s/%s/hooks/git/pre-receive", owner.Name, repo.Name).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusNoContent)

		req = NewRequestf(t, "GET", "/api/v1/repos/%s/%s/hooks/git/pre-receive", owner.Name, repo.Name).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)
		var apiGitHook2 *api.GitHook
		DecodeJSON(t, resp, &apiGitHook2)
		assert.False(t, apiGitHook2.IsActive)
		assert.Empty(t, apiGitHook2.Content)
	})

	t.Run("DeleteGitHookNoAccess", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

		session := loginUser(t, owner.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
		req := NewRequestf(t, "DELETE", "/api/v1/repos/%s/%s/hooks/git/pre-receive", owner.Name, repo.Name).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusForbidden)
	})
}
