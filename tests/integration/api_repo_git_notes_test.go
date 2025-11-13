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
	"net/url"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"

	"github.com/stretchr/testify/assert"
)

func TestAPIReposGitNotes(t *testing.T) {
	onKmupRun(t, func(*testing.T, *url.URL) {
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		// Login as User2.
		session := loginUser(t, user.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)

		// check invalid requests
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/notes/12345", user.Name).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusNotFound)

		req = NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/notes/..", user.Name).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusUnprocessableEntity)

		// check valid request
		req = NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/notes/65f1bf27bc3bf70f64657658635e66094edbcb4d", user.Name).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)

		var apiData api.Note
		DecodeJSON(t, resp, &apiData)
		assert.Equal(t, "This is a test note\n", apiData.Message)
		assert.NotEmpty(t, apiData.Commit.Files)
		assert.NotNil(t, apiData.Commit.RepoCommit.Verification)
	})
}
