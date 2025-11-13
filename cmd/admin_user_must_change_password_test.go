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
	"testing"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustChangePassword(t *testing.T) {
	defer func() {
		require.NoError(t, db.TruncateBeans(t.Context(), &user_model.User{}))
	}()
	err := microcmdUserCreate().Run(t.Context(), []string{"create", "--username", "testuser", "--email", "testuser@kmup.local", "--random-password"})
	require.NoError(t, err)
	err = microcmdUserCreate().Run(t.Context(), []string{"create", "--username", "testuserexclude", "--email", "testuserexclude@kmup.local", "--random-password"})
	require.NoError(t, err)
	// Reset password change flag
	err = microcmdUserMustChangePassword().Run(t.Context(), []string{"change-test", "--all", "--unset"})
	require.NoError(t, err)

	testUser := unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
	assert.False(t, testUser.MustChangePassword)
	testUserExclude := unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuserexclude"})
	assert.False(t, testUserExclude.MustChangePassword)

	// Make all users change password
	err = microcmdUserMustChangePassword().Run(t.Context(), []string{"change-test", "--all"})
	require.NoError(t, err)

	testUser = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
	assert.True(t, testUser.MustChangePassword)
	testUserExclude = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuserexclude"})
	assert.True(t, testUserExclude.MustChangePassword)

	// Reset password change flag but exclude all tested users
	err = microcmdUserMustChangePassword().Run(t.Context(), []string{"change-test", "--all", "--unset", "--exclude", "testuser,testuserexclude"})
	require.NoError(t, err)

	testUser = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
	assert.True(t, testUser.MustChangePassword)
	testUserExclude = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuserexclude"})
	assert.True(t, testUserExclude.MustChangePassword)

	// Reset password change flag by listing multiple users
	err = microcmdUserMustChangePassword().Run(t.Context(), []string{"change-test", "--unset", "testuser", "testuserexclude"})
	require.NoError(t, err)

	testUser = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
	assert.False(t, testUser.MustChangePassword)
	testUserExclude = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuserexclude"})
	assert.False(t, testUserExclude.MustChangePassword)

	// Exclude a user from all user
	err = microcmdUserMustChangePassword().Run(t.Context(), []string{"change-test", "--all", "--exclude", "testuserexclude"})
	require.NoError(t, err)

	testUser = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
	assert.True(t, testUser.MustChangePassword)
	testUserExclude = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuserexclude"})
	assert.False(t, testUserExclude.MustChangePassword)

	// Unset a flag for single user
	err = microcmdUserMustChangePassword().Run(t.Context(), []string{"change-test", "--unset", "testuser"})
	require.NoError(t, err)

	testUser = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
	assert.False(t, testUser.MustChangePassword)
	testUserExclude = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuserexclude"})
	assert.False(t, testUserExclude.MustChangePassword)
}
