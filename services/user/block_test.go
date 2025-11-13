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

package user

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestCanBlockUser(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
	user29 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 29})
	org3 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3})

	// Doer can't self block
	assert.False(t, CanBlockUser(t.Context(), user1, user2, user1))
	// Blocker can't be blockee
	assert.False(t, CanBlockUser(t.Context(), user1, user2, user2))
	// Can't block already blocked user
	assert.False(t, CanBlockUser(t.Context(), user1, user2, user29))
	// Blockee can't be an organization
	assert.False(t, CanBlockUser(t.Context(), user1, user2, org3))
	// Doer must be blocker or admin
	assert.False(t, CanBlockUser(t.Context(), user2, user4, user29))
	// Organization can't block a member
	assert.False(t, CanBlockUser(t.Context(), user1, org3, user4))
	// Doer must be organization owner or admin if blocker is an organization
	assert.False(t, CanBlockUser(t.Context(), user4, org3, user2))

	assert.True(t, CanBlockUser(t.Context(), user1, user2, user4))
	assert.True(t, CanBlockUser(t.Context(), user2, user2, user4))
	assert.True(t, CanBlockUser(t.Context(), user2, org3, user29))
}

func TestCanUnblockUser(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	user28 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 28})
	user29 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 29})
	org17 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 17})

	// Doer can't self unblock
	assert.False(t, CanUnblockUser(t.Context(), user1, user2, user1))
	// Can't unblock not blocked user
	assert.False(t, CanUnblockUser(t.Context(), user1, user2, user28))
	// Doer must be blocker or admin
	assert.False(t, CanUnblockUser(t.Context(), user28, user2, user29))
	// Doer must be organization owner or admin if blocker is an organization
	assert.False(t, CanUnblockUser(t.Context(), user2, org17, user28))

	assert.True(t, CanUnblockUser(t.Context(), user1, user2, user29))
	assert.True(t, CanUnblockUser(t.Context(), user2, user2, user29))
	assert.True(t, CanUnblockUser(t.Context(), user1, org17, user28))
}
