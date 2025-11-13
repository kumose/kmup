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
	"time"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/stretchr/testify/assert"
)

func TestCheckAuthToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	t.Run("Empty", func(t *testing.T) {
		token, err := CheckAuthToken(t.Context(), "")
		assert.NoError(t, err)
		assert.Nil(t, token)
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		token, err := CheckAuthToken(t.Context(), "dummy")
		assert.ErrorIs(t, err, ErrAuthTokenInvalidFormat)
		assert.Nil(t, token)
	})

	t.Run("NotFound", func(t *testing.T) {
		token, err := CheckAuthToken(t.Context(), "notexists:dummy")
		assert.ErrorIs(t, err, ErrAuthTokenExpired)
		assert.Nil(t, token)
	})

	t.Run("Expired", func(t *testing.T) {
		timeutil.MockSet(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

		at, token, err := CreateAuthTokenForUserID(t.Context(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, at)
		assert.NotEmpty(t, token)

		timeutil.MockUnset()

		at2, err := CheckAuthToken(t.Context(), at.ID+":"+token)
		assert.ErrorIs(t, err, ErrAuthTokenExpired)
		assert.Nil(t, at2)

		assert.NoError(t, auth_model.DeleteAuthTokenByID(t.Context(), at.ID))
	})

	t.Run("InvalidHash", func(t *testing.T) {
		at, token, err := CreateAuthTokenForUserID(t.Context(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, at)
		assert.NotEmpty(t, token)

		at2, err := CheckAuthToken(t.Context(), at.ID+":"+token+"dummy")
		assert.ErrorIs(t, err, ErrAuthTokenInvalidHash)
		assert.Nil(t, at2)

		assert.NoError(t, auth_model.DeleteAuthTokenByID(t.Context(), at.ID))
	})

	t.Run("Valid", func(t *testing.T) {
		at, token, err := CreateAuthTokenForUserID(t.Context(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, at)
		assert.NotEmpty(t, token)

		at2, err := CheckAuthToken(t.Context(), at.ID+":"+token)
		assert.NoError(t, err)
		assert.NotNil(t, at2)

		assert.NoError(t, auth_model.DeleteAuthTokenByID(t.Context(), at.ID))
	})
}

func TestRegenerateAuthToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	timeutil.MockSet(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	defer timeutil.MockUnset()

	at, token, err := CreateAuthTokenForUserID(t.Context(), 2)
	assert.NoError(t, err)
	assert.NotNil(t, at)
	assert.NotEmpty(t, token)

	timeutil.MockSet(time.Date(2023, 1, 1, 0, 0, 1, 0, time.UTC))

	at2, token2, err := RegenerateAuthToken(t.Context(), at)
	assert.NoError(t, err)
	assert.NotNil(t, at2)
	assert.NotEmpty(t, token2)

	assert.Equal(t, at.ID, at2.ID)
	assert.Equal(t, at.UserID, at2.UserID)
	assert.NotEqual(t, token, token2)
	assert.NotEqual(t, at.ExpiresUnix, at2.ExpiresUnix)

	assert.NoError(t, auth_model.DeleteAuthTokenByID(t.Context(), at.ID))
}
