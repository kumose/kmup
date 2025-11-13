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

package v1_22

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/models/project"

	"github.com/stretchr/testify/assert"
)

func Test_CheckProjectColumnsConsistency(t *testing.T) {
	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(project.Project), new(project.Column))
	defer deferable()
	if x == nil || t.Failed() {
		return
	}

	assert.NoError(t, CheckProjectColumnsConsistency(x))

	// check if default column was added
	var defaultColumn project.Column
	has, err := x.Where("project_id=? AND `default` = ?", 1, true).Get(&defaultColumn)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.Equal(t, int64(1), defaultColumn.ProjectID)
	assert.True(t, defaultColumn.Default)

	// check if multiple defaults, previous were removed and last will be kept
	expectDefaultColumn, err := project.GetColumn(t.Context(), 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), expectDefaultColumn.ProjectID)
	assert.False(t, expectDefaultColumn.Default)

	expectNonDefaultColumn, err := project.GetColumn(t.Context(), 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), expectNonDefaultColumn.ProjectID)
	assert.True(t, expectNonDefaultColumn.Default)
}
