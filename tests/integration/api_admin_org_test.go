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
	"strings"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIAdminOrgCreate(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	session := loginUser(t, "user1")
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteAdmin)

	t.Run("CreateOrg", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		org := api.CreateOrgOption{
			UserName:    "user2_org",
			FullName:    "User2's organization",
			Description: "This organization created by admin for user2",
			Website:     "https://try.kmup.io",
			Location:    "Shanghai",
			Visibility:  "private",
		}
		req := NewRequestWithJSON(t, "POST", "/api/v1/admin/users/user2/orgs", &org).
			AddTokenAuth(token)
		resp := MakeRequest(t, req, http.StatusCreated)

		var apiOrg api.Organization
		DecodeJSON(t, resp, &apiOrg)

		assert.Equal(t, org.UserName, apiOrg.Name)
		assert.Equal(t, org.FullName, apiOrg.FullName)
		assert.Equal(t, org.Description, apiOrg.Description)
		assert.Equal(t, org.Website, apiOrg.Website)
		assert.Equal(t, org.Location, apiOrg.Location)
		assert.Equal(t, org.Visibility, apiOrg.Visibility)

		unittest.AssertExistsAndLoadBean(t, &user_model.User{
			Name:      org.UserName,
			LowerName: strings.ToLower(org.UserName),
			FullName:  org.FullName,
		})
	})
	t.Run("CreateBadVisibility", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		org := api.CreateOrgOption{
			UserName:    "user2_org",
			FullName:    "User2's organization",
			Description: "This organization created by admin for user2",
			Website:     "https://try.kmup.io",
			Location:    "Shanghai",
			Visibility:  "notvalid",
		}
		req := NewRequestWithJSON(t, "POST", "/api/v1/admin/users/user2/orgs", &org).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusUnprocessableEntity)
	})
	t.Run("CreateNotAdmin", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		nonAdminUsername := "user2"
		session := loginUser(t, nonAdminUsername)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeAll)
		org := api.CreateOrgOption{
			UserName:    "user2_org",
			FullName:    "User2's organization",
			Description: "This organization created by admin for user2",
			Website:     "https://try.kmup.io",
			Location:    "Shanghai",
			Visibility:  "public",
		}
		req := NewRequestWithJSON(t, "POST", "/api/v1/admin/users/user2/orgs", &org).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusForbidden)
	})
}
