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

package auth

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/reqctx"
	"github.com/kumose/kmup/services/actions"

	"github.com/stretchr/testify/assert"
)

func TestUserIDFromToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	t.Run("Actions JWT", func(t *testing.T) {
		const RunningTaskID = 47
		token, err := actions.CreateAuthorizationToken(RunningTaskID, 1, 2)
		assert.NoError(t, err)

		ds := make(reqctx.ContextData)

		o := OAuth2{}
		uid := o.userIDFromToken(t.Context(), token, ds)
		assert.Equal(t, user_model.ActionsUserID, uid)
		assert.Equal(t, true, ds["IsActionsToken"])
		assert.Equal(t, ds["ActionsTaskID"], int64(RunningTaskID))
	})
}

func TestCheckTaskIsRunning(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	cases := map[string]struct {
		TaskID   int64
		Expected bool
	}{
		"Running":   {TaskID: 47, Expected: true},
		"Missing":   {TaskID: 1, Expected: false},
		"Cancelled": {TaskID: 46, Expected: false},
	}

	for name := range cases {
		c := cases[name]
		t.Run(name, func(t *testing.T) {
			actual := CheckTaskIsRunning(t.Context(), c.TaskID)
			assert.Equal(t, c.Expected, actual)
		})
	}
}
