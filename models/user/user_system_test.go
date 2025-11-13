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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemUser(t *testing.T) {
	u, err := GetPossibleUserByID(t.Context(), -1)
	require.NoError(t, err)
	assert.Equal(t, "Ghost", u.Name)
	assert.Equal(t, "ghost", u.LowerName)
	assert.True(t, u.IsGhost())
	assert.True(t, IsGhostUserName("gHost"))

	u, err = GetPossibleUserByID(t.Context(), -2)
	require.NoError(t, err)
	assert.Equal(t, "kmup-actions", u.Name)
	assert.Equal(t, "kmup-actions", u.LowerName)
	assert.True(t, u.IsKmupActions())
	assert.True(t, IsKmupActionsUserName("Kmup-actionS"))

	_, err = GetPossibleUserByID(t.Context(), -3)
	require.Error(t, err)
}
