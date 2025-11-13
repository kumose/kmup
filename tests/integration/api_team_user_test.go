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
	"time"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/services/convert"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPITeamUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	user2Session := loginUser(t, "user2")
	user2Token := getTokenForLoggedInUser(t, user2Session, auth_model.AccessTokenScopeWriteOrganization)

	t.Run("User2ReadUser1", func(t *testing.T) {
		req := NewRequest(t, "GET", "/api/v1/teams/1/members/user1").AddTokenAuth(user2Token)
		MakeRequest(t, req, http.StatusNotFound)
	})

	t.Run("User2ReadSelf", func(t *testing.T) {
		// read self user
		req := NewRequest(t, "GET", "/api/v1/teams/1/members/user2").AddTokenAuth(user2Token)
		resp := MakeRequest(t, req, http.StatusOK)
		var user2 *api.User
		DecodeJSON(t, resp, &user2)
		user2.Created = user2.Created.In(time.Local)
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: "user2"})

		expectedUser := convert.ToUser(t.Context(), user, user)

		// test time via unix timestamp
		assert.Equal(t, expectedUser.LastLogin.Unix(), user2.LastLogin.Unix())
		assert.Equal(t, expectedUser.Created.Unix(), user2.Created.Unix())
		expectedUser.LastLogin = user2.LastLogin
		expectedUser.Created = user2.Created

		assert.Equal(t, expectedUser, user2)
	})
}
