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

func TestChangePasswordCommand(t *testing.T) {
	ctx := t.Context()

	defer func() {
		require.NoError(t, db.TruncateBeans(t.Context(), &user_model.User{}))
	}()

	t.Run("change password successfully", func(t *testing.T) {
		// defer func() {
		// 	require.NoError(t, db.TruncateBeans(t.Context(), &user_model.User{}))
		// }()
		// Prepare test user
		unittest.AssertNotExistsBean(t, &user_model.User{LowerName: "testuser"})
		err := microcmdUserCreate().Run(ctx, []string{"create", "--username", "testuser", "--email", "testuser@kmup.local", "--random-password"})
		require.NoError(t, err)

		// load test user
		userBase := unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})

		// Change the password
		err = microcmdUserChangePassword().Run(ctx, []string{"change-password", "--username", "testuser", "--password", "newpassword"})
		require.NoError(t, err)

		// Verify the password has been changed
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
		assert.NotEqual(t, userBase.Passwd, user.Passwd)
		assert.NotEqual(t, userBase.Salt, user.Salt)

		// Additional check for must-change-password flag
		require.NoError(t, microcmdUserChangePassword().Run(ctx, []string{"change-password", "--username", "testuser", "--password", "anotherpassword", "--must-change-password=false"}))
		user = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
		assert.False(t, user.MustChangePassword)

		require.NoError(t, microcmdUserChangePassword().Run(ctx, []string{"change-password", "--username", "testuser", "--password", "yetanotherpassword", "--must-change-password"}))
		user = unittest.AssertExistsAndLoadBean(t, &user_model.User{LowerName: "testuser"})
		assert.True(t, user.MustChangePassword)
	})

	t.Run("failure cases", func(t *testing.T) {
		testCases := []struct {
			name        string
			args        []string
			expectedErr string
		}{
			{
				name:        "user does not exist",
				args:        []string{"change-password", "--username", "nonexistentuser", "--password", "newpassword"},
				expectedErr: "user does not exist",
			},
			{
				name:        "missing username",
				args:        []string{"change-password", "--password", "newpassword"},
				expectedErr: `"username" not set`,
			},
			{
				name:        "missing password",
				args:        []string{"change-password", "--username", "testuser"},
				expectedErr: `"password" not set`,
			},
			{
				name:        "too short password",
				args:        []string{"change-password", "--username", "testuser", "--password", "1"},
				expectedErr: "password is not long enough",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := microcmdUserChangePassword().Run(ctx, tc.args)
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			})
		}
	})
}
