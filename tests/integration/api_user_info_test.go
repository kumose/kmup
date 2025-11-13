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

func TestAPIUserInfo(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := "user1"
	user2 := "user31"

	org3 := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: "org3"})

	session := loginUser(t, user)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadUser)

	t.Run("GetInfo", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequest(t, "GET", "/api/v1/users/"+user2).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)

		var u api.User
		DecodeJSON(t, resp, &u)
		assert.Equal(t, user2, u.UserName)

		req = NewRequest(t, "GET", "/api/v1/users/"+user2)
		MakeRequest(t, req, http.StatusNotFound)

		// test if the placaholder Mail is returned if a User is not logged in
		req = NewRequest(t, "GET", "/api/v1/users/"+org3.Name)
		resp = MakeRequest(t, req, http.StatusOK)
		DecodeJSON(t, resp, &u)
		assert.Equal(t, org3.GetPlaceholderEmail(), u.Email)

		// Test if the correct Mail is returned if a User is logged in
		req = NewRequest(t, "GET", "/api/v1/users/"+org3.Name).
			AddTokenAuth(token)
		resp = MakeRequest(t, req, http.StatusOK)
		DecodeJSON(t, resp, &u)
		assert.Equal(t, org3.GetEmail(), u.Email)
	})

	t.Run("GetAuthenticatedUser", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequest(t, "GET", "/api/v1/user").
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusOK)

		var u api.User
		DecodeJSON(t, resp, &u)
		assert.Equal(t, user, u.UserName)
	})
}
