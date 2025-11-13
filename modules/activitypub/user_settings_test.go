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

package activitypub

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	_ "github.com/kumose/kmup/models"

	"github.com/stretchr/testify/assert"
)

func TestUserSettings(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	pub, priv, err := GetKeyPair(t.Context(), user1)
	assert.NoError(t, err)
	pub1, err := GetPublicKey(t.Context(), user1)
	assert.NoError(t, err)
	assert.Equal(t, pub, pub1)
	priv1, err := GetPrivateKey(t.Context(), user1)
	assert.NoError(t, err)
	assert.Equal(t, priv, priv1)
}
