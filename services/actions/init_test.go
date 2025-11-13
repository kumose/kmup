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

package actions

import (
	"os"
	"testing"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
	os.Exit(m.Run())
}

func TestInitToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	t.Run("NoToken", func(t *testing.T) {
		_, _ = db.Exec(t.Context(), "DELETE FROM action_runner_token")
		t.Setenv("KMUP_RUNNER_REGISTRATION_TOKEN", "")
		t.Setenv("KMUP_RUNNER_REGISTRATION_TOKEN_FILE", "")
		err := initGlobalRunnerToken(t.Context())
		require.NoError(t, err)
		notEmpty, err := db.IsTableNotEmpty(&actions_model.ActionRunnerToken{})
		require.NoError(t, err)
		assert.False(t, notEmpty)
	})

	t.Run("EnvToken", func(t *testing.T) {
		tokenValue, _ := util.CryptoRandomString(32)
		t.Setenv("KMUP_RUNNER_REGISTRATION_TOKEN", tokenValue)
		t.Setenv("KMUP_RUNNER_REGISTRATION_TOKEN_FILE", "")
		err := initGlobalRunnerToken(t.Context())
		require.NoError(t, err)
		token := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunnerToken{Token: tokenValue})
		assert.True(t, token.IsActive)

		// init with the same token again, should not create a new token
		err = initGlobalRunnerToken(t.Context())
		require.NoError(t, err)
		token2 := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunnerToken{Token: tokenValue})
		assert.Equal(t, token.ID, token2.ID)
		assert.True(t, token.IsActive)
	})

	t.Run("EnvFileToken", func(t *testing.T) {
		tokenValue, _ := util.CryptoRandomString(32)
		f := t.TempDir() + "/token"
		_ = os.WriteFile(f, []byte(tokenValue), 0o644)
		t.Setenv("KMUP_RUNNER_REGISTRATION_TOKEN", "")
		t.Setenv("KMUP_RUNNER_REGISTRATION_TOKEN_FILE", f)
		err := initGlobalRunnerToken(t.Context())
		require.NoError(t, err)
		token := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunnerToken{Token: tokenValue})
		assert.True(t, token.IsActive)

		// if the env token is invalidated by another new token, then it shouldn't be active anymore
		_, err = actions_model.NewRunnerToken(t.Context(), 0, 0)
		require.NoError(t, err)
		token = unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunnerToken{Token: tokenValue})
		assert.False(t, token.IsActive)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		t.Setenv("KMUP_RUNNER_REGISTRATION_TOKEN", "abc")
		err := initGlobalRunnerToken(t.Context())
		assert.ErrorContains(t, err, "must be at least")
	})
}
