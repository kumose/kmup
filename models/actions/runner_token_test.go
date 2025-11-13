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
	"testing"

	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestRunnerToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	token := unittest.AssertExistsAndLoadBean(t, &ActionRunnerToken{ID: 3})
	expectedToken, err := GetLatestRunnerToken(t.Context(), 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestNewRunnerToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	token, err := NewRunnerToken(t.Context(), 1, 0)
	assert.NoError(t, err)
	expectedToken, err := GetLatestRunnerToken(t.Context(), 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestUpdateRunnerToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	token := unittest.AssertExistsAndLoadBean(t, &ActionRunnerToken{ID: 3})
	token.IsActive = true
	assert.NoError(t, UpdateRunnerToken(t.Context(), token))
	expectedToken, err := GetLatestRunnerToken(t.Context(), 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}
