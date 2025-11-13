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
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/tests"
)

func TestAPIReposGitRefs(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)

	for _, ref := range [...]string{
		"refs/heads/master", // Branch
		"refs/tags/v1.1",    // Tag
	} {
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/%s", user.Name, ref).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusOK)
	}
	// Test getting all refs
	req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/refs", user.Name).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusOK)
	// Test getting non-existent refs
	req = NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/refs/heads/unknown", user.Name).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNotFound)
}
