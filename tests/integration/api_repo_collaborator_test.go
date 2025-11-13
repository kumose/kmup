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
	"github.com/kumose/kmup/models/perm"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIRepoCollaboratorPermission(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	repo2 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	repo2Owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo2.OwnerID})

	user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
	user5 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 5})
	user10 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 10})
	user11 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 11})
	user34 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 34})

	testCtx := NewAPITestContext(t, repo2Owner.Name, repo2.Name, auth_model.AccessTokenScopeWriteRepository)

	t.Run("RepoOwnerShouldBeOwner", func(t *testing.T) {
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, repo2Owner.Name).
			AddTokenAuth(testCtx.Token)
		resp := MakeRequest(t, req, http.StatusOK)

		var repoPermission api.RepoCollaboratorPermission
		DecodeJSON(t, resp, &repoPermission)

		assert.Equal(t, "owner", repoPermission.Permission)
	})

	t.Run("CollaboratorWithReadAccess", func(t *testing.T) {
		t.Run("AddUserAsCollaboratorWithReadAccess", doAPIAddCollaborator(testCtx, user4.Name, perm.AccessModeRead))

		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, user4.Name).
			AddTokenAuth(testCtx.Token)
		resp := MakeRequest(t, req, http.StatusOK)

		var repoPermission api.RepoCollaboratorPermission
		DecodeJSON(t, resp, &repoPermission)

		assert.Equal(t, "read", repoPermission.Permission)
	})

	t.Run("CollaboratorWithWriteAccess", func(t *testing.T) {
		t.Run("AddUserAsCollaboratorWithWriteAccess", doAPIAddCollaborator(testCtx, user4.Name, perm.AccessModeWrite))

		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, user4.Name).
			AddTokenAuth(testCtx.Token)
		resp := MakeRequest(t, req, http.StatusOK)

		var repoPermission api.RepoCollaboratorPermission
		DecodeJSON(t, resp, &repoPermission)

		assert.Equal(t, "write", repoPermission.Permission)
	})

	t.Run("CollaboratorWithAdminAccess", func(t *testing.T) {
		t.Run("AddUserAsCollaboratorWithAdminAccess", doAPIAddCollaborator(testCtx, user4.Name, perm.AccessModeAdmin))

		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, user4.Name).
			AddTokenAuth(testCtx.Token)
		resp := MakeRequest(t, req, http.StatusOK)

		var repoPermission api.RepoCollaboratorPermission
		DecodeJSON(t, resp, &repoPermission)

		assert.Equal(t, "admin", repoPermission.Permission)
	})

	t.Run("CollaboratorNotFound", func(t *testing.T) {
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, "non-existent-user").
			AddTokenAuth(testCtx.Token)
		MakeRequest(t, req, http.StatusNotFound)
	})

	t.Run("CollaboratorBlocked", func(t *testing.T) {
		ctx := NewAPITestContext(t, repo2Owner.Name, repo2.Name, auth_model.AccessTokenScopeWriteRepository)
		ctx.ExpectedCode = http.StatusForbidden
		doAPIAddCollaborator(ctx, user34.Name, perm.AccessModeAdmin)(t)
	})

	t.Run("CollaboratorCanQueryItsPermissions", func(t *testing.T) {
		t.Run("AddUserAsCollaboratorWithReadAccess", doAPIAddCollaborator(testCtx, user5.Name, perm.AccessModeRead))

		_session := loginUser(t, user5.Name)
		_testCtx := NewAPITestContext(t, user5.Name, repo2.Name, auth_model.AccessTokenScopeReadRepository)

		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, user5.Name).
			AddTokenAuth(_testCtx.Token)
		resp := _session.MakeRequest(t, req, http.StatusOK)

		var repoPermission api.RepoCollaboratorPermission
		DecodeJSON(t, resp, &repoPermission)

		assert.Equal(t, "read", repoPermission.Permission)

		t.Run("CollaboratorCanReadOwnPermission", func(t *testing.T) {
			session := loginUser(t, user5.Name)
			token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)

			req = NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, user5.Name).AddTokenAuth(token)
			resp = MakeRequest(t, req, http.StatusOK)

			repoCollPerm := api.RepoCollaboratorPermission{}
			DecodeJSON(t, resp, &repoCollPerm)

			assert.Equal(t, "read", repoCollPerm.Permission)
		})
	})

	t.Run("CollaboratorCanQueryItsPermissions", func(t *testing.T) {
		t.Run("AddUserAsCollaboratorWithReadAccess", doAPIAddCollaborator(testCtx, user5.Name, perm.AccessModeRead))

		_session := loginUser(t, user5.Name)
		_testCtx := NewAPITestContext(t, user5.Name, repo2.Name, auth_model.AccessTokenScopeReadRepository)

		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, user5.Name).
			AddTokenAuth(_testCtx.Token)
		resp := _session.MakeRequest(t, req, http.StatusOK)

		var repoPermission api.RepoCollaboratorPermission
		DecodeJSON(t, resp, &repoPermission)

		assert.Equal(t, "read", repoPermission.Permission)
	})

	t.Run("RepoAdminCanQueryACollaboratorsPermissions", func(t *testing.T) {
		t.Run("AddUserAsCollaboratorWithAdminAccess", doAPIAddCollaborator(testCtx, user10.Name, perm.AccessModeAdmin))
		t.Run("AddUserAsCollaboratorWithReadAccess", doAPIAddCollaborator(testCtx, user11.Name, perm.AccessModeRead))

		_session := loginUser(t, user10.Name)
		_testCtx := NewAPITestContext(t, user10.Name, repo2.Name, auth_model.AccessTokenScopeReadRepository)

		req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/collaborators/%s/permission", repo2Owner.Name, repo2.Name, user11.Name).
			AddTokenAuth(_testCtx.Token)
		resp := _session.MakeRequest(t, req, http.StatusOK)

		var repoPermission api.RepoCollaboratorPermission
		DecodeJSON(t, resp, &repoPermission)

		assert.Equal(t, "read", repoPermission.Permission)
	})
}
