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
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPICompareBranches(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

	t.Run("CompareBranches", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		req := NewRequestf(t, "GET", "/api/v1/repos/user2/repo20/compare/add-csv...remove-files-b").AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)

		var apiResp *api.Compare
		DecodeJSON(t, resp, &apiResp)

		assert.Equal(t, 2, apiResp.TotalCommits)
		assert.Len(t, apiResp.Commits, 2)
	})

	t.Run("CompareCommits", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		req := NewRequestf(t, "GET", "/api/v1/repos/user2/repo20/compare/808038d2f71b0ab02099...c8e31bc7688741a5287f").AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)

		var apiResp *api.Compare
		DecodeJSON(t, resp, &apiResp)

		assert.Equal(t, 1, apiResp.TotalCommits)
		assert.Len(t, apiResp.Commits, 1)
	})
}
