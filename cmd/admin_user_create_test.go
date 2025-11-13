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

package cmd

import (
	"fmt"
	"strings"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminUserCreate(t *testing.T) {
	reset := func() {
		require.NoError(t, db.TruncateBeans(t.Context(), &user_model.User{}))
		require.NoError(t, db.TruncateBeans(t.Context(), &user_model.EmailAddress{}))
		require.NoError(t, db.TruncateBeans(t.Context(), &auth_model.AccessToken{}))
	}

	t.Run("MustChangePassword", func(t *testing.T) {
		type check struct {
			IsAdmin            bool
			MustChangePassword bool
		}

		createCheck := func(name, args string) check {
			require.NoError(t, microcmdUserCreate().Run(t.Context(), strings.Fields(fmt.Sprintf("create --username %s --email %s@kmup.local %s --password foobar", name, name, args))))
			u := unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: name})
			return check{IsAdmin: u.IsAdmin, MustChangePassword: u.MustChangePassword}
		}
		reset()
		assert.Equal(t, check{IsAdmin: false, MustChangePassword: false}, createCheck("u", ""), "first non-admin user doesn't need to change password")

		reset()
		assert.Equal(t, check{IsAdmin: true, MustChangePassword: false}, createCheck("u", "--admin"), "first admin user doesn't need to change password")

		reset()
		assert.Equal(t, check{IsAdmin: true, MustChangePassword: true}, createCheck("u", "--admin --must-change-password"))
		assert.Equal(t, check{IsAdmin: true, MustChangePassword: true}, createCheck("u2", "--admin"))
		assert.Equal(t, check{IsAdmin: true, MustChangePassword: false}, createCheck("u3", "--admin --must-change-password=false"))
		assert.Equal(t, check{IsAdmin: false, MustChangePassword: true}, createCheck("u4", ""))
		assert.Equal(t, check{IsAdmin: false, MustChangePassword: false}, createCheck("u5", "--must-change-password=false"))
	})

	createUser := func(name string, args ...string) error {
		return microcmdUserCreate().Run(t.Context(), append([]string{"create", "--username", name, "--email", name + "@kmup.local"}, args...))
	}

	t.Run("UserType", func(t *testing.T) {
		reset()
		assert.ErrorContains(t, createUser("u", "--user-type", "invalid"), "invalid user type")
		assert.ErrorContains(t, createUser("u", "--user-type", "bot", "--password", "123"), "can only be set for individual users")
		assert.ErrorContains(t, createUser("u", "--user-type", "bot", "--must-change-password"), "can only be set for individual users")

		assert.NoError(t, createUser("u", "--user-type", "bot"))
		u := unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "u"})
		assert.Equal(t, user_model.UserTypeBot, u.Type)
		assert.Empty(t, u.Passwd)
	})

	t.Run("AccessToken", func(t *testing.T) {
		// no generated access token
		reset()
		assert.NoError(t, createUser("u", "--random-password"))
		assert.Equal(t, 1, unittest.GetCount(t, &user_model.User{}))
		assert.Equal(t, 0, unittest.GetCount(t, &auth_model.AccessToken{}))

		// using "--access-token" only means "all" access
		reset()
		assert.NoError(t, createUser("u", "--random-password", "--access-token"))
		assert.Equal(t, 1, unittest.GetCount(t, &user_model.User{}))
		assert.Equal(t, 1, unittest.GetCount(t, &auth_model.AccessToken{}))
		accessToken := unittest.AssertExistsAndLoadBean(t, &auth_model.AccessToken{Name: "kmup-admin"})
		hasScopes, err := accessToken.Scope.HasScope(auth_model.AccessTokenScopeWriteAdmin, auth_model.AccessTokenScopeWriteRepository)
		assert.NoError(t, err)
		assert.True(t, hasScopes)

		// using "--access-token" with name & scopes
		reset()
		assert.NoError(t, createUser("u", "--random-password", "--access-token", "--access-token-name", "new-token-name", "--access-token-scopes", "read:issue,read:user"))
		assert.Equal(t, 1, unittest.GetCount(t, &user_model.User{}))
		assert.Equal(t, 1, unittest.GetCount(t, &auth_model.AccessToken{}))
		accessToken = unittest.AssertExistsAndLoadBean(t, &auth_model.AccessToken{Name: "new-token-name"})
		hasScopes, err = accessToken.Scope.HasScope(auth_model.AccessTokenScopeReadIssue, auth_model.AccessTokenScopeReadUser)
		assert.NoError(t, err)
		assert.True(t, hasScopes)
		hasScopes, err = accessToken.Scope.HasScope(auth_model.AccessTokenScopeWriteAdmin, auth_model.AccessTokenScopeWriteRepository)
		assert.NoError(t, err)
		assert.False(t, hasScopes)

		// using "--access-token-name" without "--access-token"
		reset()
		err = createUser("u", "--random-password", "--access-token-name", "new-token-name")
		assert.Equal(t, 0, unittest.GetCount(t, &user_model.User{}))
		assert.Equal(t, 0, unittest.GetCount(t, &auth_model.AccessToken{}))
		assert.ErrorContains(t, err, "access-token-name and access-token-scopes flags are only valid when access-token flag is set")

		// using "--access-token-scopes" without "--access-token"
		reset()
		err = createUser("u", "--random-password", "--access-token-scopes", "read:issue")
		assert.Equal(t, 0, unittest.GetCount(t, &user_model.User{}))
		assert.Equal(t, 0, unittest.GetCount(t, &auth_model.AccessToken{}))
		assert.ErrorContains(t, err, "access-token-name and access-token-scopes flags are only valid when access-token flag is set")

		// empty permission
		reset()
		err = createUser("u", "--random-password", "--access-token", "--access-token-scopes", "public-only")
		assert.Equal(t, 0, unittest.GetCount(t, &user_model.User{}))
		assert.Equal(t, 0, unittest.GetCount(t, &auth_model.AccessToken{}))
		assert.ErrorContains(t, err, "access token does not have any permission")
	})

	t.Run("UserFields", func(t *testing.T) {
		reset()
		assert.NoError(t, createUser("u-FullNameWithSpace", "--random-password", "--fullname", "First O'Middle Last"))
		unittest.AssertExistsAndLoadBean(t, &user_model.User{
			Name:      "u-FullNameWithSpace",
			LowerName: "u-fullnamewithspace",
			FullName:  "First O'Middle Last",
			Email:     "u-FullNameWithSpace@kmup.local",
		})

		assert.NoError(t, createUser("u-FullNameEmpty", "--random-password", "--fullname", ""))
		u := unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "u-fullnameempty"})
		assert.Empty(t, u.FullName)
	})
}
