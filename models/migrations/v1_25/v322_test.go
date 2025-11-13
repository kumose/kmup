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

package v1_25

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ExtendCommentTreePathLength(t *testing.T) {
	if setting.Database.Type.IsSQLite3() {
		t.Skip("For SQLITE, varchar or char will always be represented as TEXT")
	}

	type Comment struct {
		ID       int64  `xorm:"pk autoincr"`
		TreePath string `xorm:"VARCHAR(255)"`
	}

	x, deferrable := base.PrepareTestEnv(t, 0, new(Comment))
	defer deferrable()

	require.NoError(t, ExtendCommentTreePathLength(x))
	table := base.LoadTableSchemasMap(t, x)["comment"]
	column := table.GetColumn("tree_path")
	assert.Contains(t, []string{"NVARCHAR", "VARCHAR"}, column.SQLType.Name)
	assert.EqualValues(t, 4000, column.Length)
}
