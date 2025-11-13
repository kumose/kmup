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

package migrations

import (
	"testing"

	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestMigrations(t *testing.T) {
	defer test.MockVariableValue(&preparedMigrations)()
	preparedMigrations = []*migration{
		{idNumber: 70},
		{idNumber: 71},
	}
	assert.EqualValues(t, 72, calcDBVersion(preparedMigrations))
	assert.EqualValues(t, 72, ExpectedDBVersion())

	assert.EqualValues(t, 71, migrationIDNumberToDBVersion(70))

	assert.Equal(t, []*migration{{idNumber: 70}, {idNumber: 71}}, getPendingMigrations(70, preparedMigrations))
	assert.Equal(t, []*migration{{idNumber: 71}}, getPendingMigrations(71, preparedMigrations))
	assert.Equal(t, []*migration{}, getPendingMigrations(72, preparedMigrations))
}
