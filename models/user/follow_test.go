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

package user_test

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestIsFollowing(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.True(t, user_model.IsFollowing(t.Context(), 4, 2))
	assert.False(t, user_model.IsFollowing(t.Context(), 2, 4))
	assert.False(t, user_model.IsFollowing(t.Context(), 5, unittest.NonexistentID))
	assert.False(t, user_model.IsFollowing(t.Context(), unittest.NonexistentID, 5))
	assert.False(t, user_model.IsFollowing(t.Context(), unittest.NonexistentID, unittest.NonexistentID))
}
