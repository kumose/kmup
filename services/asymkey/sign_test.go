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

package asymkey

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHasPubkeys(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	test := func(t *testing.T, userID int64, expectedHasGPG, expectedHasSSH bool) {
		ctx := t.Context()
		hasGPG, err := userHasPubkeysGPG(ctx, userID)
		require.NoError(t, err)
		hasSSH, err := userHasPubkeysSSH(ctx, userID)
		require.NoError(t, err)
		hasPubkeys, err := userHasPubkeys(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, expectedHasGPG, hasGPG)
		assert.Equal(t, expectedHasSSH, hasSSH)
		assert.Equal(t, expectedHasGPG || expectedHasSSH, hasPubkeys)
	}

	t.Run("AllowUserWithGPGKey", func(t *testing.T) {
		test(t, 36, true, false) // has gpg
	})
	t.Run("AllowUserWithSSHKey", func(t *testing.T) {
		test(t, 2, false, true) // has ssh
	})
	t.Run("DenyUserWithNoKeys", func(t *testing.T) {
		test(t, 1, false, false) // no pubkey
	})
}
