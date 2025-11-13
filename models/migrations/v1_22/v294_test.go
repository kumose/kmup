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

	"github.com/stretchr/testify/assert"
	"xorm.io/xorm/schemas"
)

func Test_AddUniqueIndexForProjectIssue(t *testing.T) {
	type ProjectIssue struct { //revive:disable-line:exported
		ID        int64 `xorm:"pk autoincr"`
		IssueID   int64 `xorm:"INDEX"`
		ProjectID int64 `xorm:"INDEX"`
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(ProjectIssue))
	defer deferable()
	if x == nil || t.Failed() {
		return
	}

	cnt, err := x.Table("project_issue").Where("project_id=1 AND issue_id=1").Count()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)

	assert.NoError(t, AddUniqueIndexForProjectIssue(x))

	cnt, err = x.Table("project_issue").Where("project_id=1 AND issue_id=1").Count()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	tables, err := x.DBMetas()
	assert.NoError(t, err)
	assert.Len(t, tables, 1)
	found := false
	for _, index := range tables[0].Indexes {
		if index.Type == schemas.UniqueType {
			found = true
			assert.ElementsMatch(t, index.Cols, []string{"project_id", "issue_id"})
			break
		}
	}
	assert.True(t, found)
}
