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

package system_test

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/system"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	keyName := "test.key"
	assert.NoError(t, unittest.PrepareTestDatabase())

	assert.NoError(t, db.TruncateBeans(t.Context(), &system.Setting{}))

	rev, settings, err := system.GetAllSettings(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, rev)
	assert.Len(t, settings, 1) // there is only one "revision" key

	err = system.SetSettings(t.Context(), map[string]string{keyName: "true"})
	assert.NoError(t, err)
	rev, settings, err = system.GetAllSettings(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 2, rev)
	assert.Len(t, settings, 2)
	assert.Equal(t, "true", settings[keyName])

	err = system.SetSettings(t.Context(), map[string]string{keyName: "false"})
	assert.NoError(t, err)
	rev, settings, err = system.GetAllSettings(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 3, rev)
	assert.Len(t, settings, 2)
	assert.Equal(t, "false", settings[keyName])

	// setting the same value should not trigger DuplicateKey error, and the "version" should be increased
	err = system.SetSettings(t.Context(), map[string]string{keyName: "false"})
	assert.NoError(t, err)

	rev, settings, err = system.GetAllSettings(t.Context())
	assert.NoError(t, err)
	assert.Len(t, settings, 2)
	assert.Equal(t, 4, rev)
}
