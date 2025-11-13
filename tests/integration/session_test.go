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
	"testing"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func Test_RegenerateSession(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	assert.NoError(t, unittest.PrepareTestDatabase())

	key := "new_key890123456"  // it must be 16 characters long
	key2 := "new_key890123457" // it must be 16 characters
	exist, err := auth.ExistSession(t.Context(), key)
	assert.NoError(t, err)
	assert.False(t, exist)

	sess, err := auth.RegenerateSession(t.Context(), "", key)
	assert.NoError(t, err)
	assert.Equal(t, key, sess.Key)
	assert.Empty(t, sess.Data)

	sess, err = auth.ReadSession(t.Context(), key2)
	assert.NoError(t, err)
	assert.Equal(t, key2, sess.Key)
	assert.Empty(t, sess.Data)
}
