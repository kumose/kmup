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

package v1_15

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

func Test_AddIssueResourceIndexTable(t *testing.T) {
	// Create the models used in the migration
	type Issue struct {
		ID     int64 `xorm:"pk autoincr"`
		RepoID int64 `xorm:"UNIQUE(s)"`
		Index  int64 `xorm:"UNIQUE(s)"`
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(Issue))
	if x == nil || t.Failed() {
		defer deferable()
		return
	}
	defer deferable()

	// Run the migration
	if err := AddIssueResourceIndexTable(x); err != nil {
		assert.NoError(t, err)
		return
	}

	type ResourceIndex struct {
		GroupID  int64 `xorm:"pk"`
		MaxIndex int64 `xorm:"index"`
	}

	start := 0
	const batchSize = 1000
	for {
		indexes := make([]ResourceIndex, 0, batchSize)
		err := x.Table("issue_index").Limit(batchSize, start).Find(&indexes)
		assert.NoError(t, err)

		for _, idx := range indexes {
			var maxIndex int
			has, err := x.SQL("SELECT max(`index`) FROM issue WHERE repo_id = ?", idx.GroupID).Get(&maxIndex)
			assert.NoError(t, err)
			assert.True(t, has)
			assert.EqualValues(t, maxIndex, idx.MaxIndex)
		}
		if len(indexes) < batchSize {
			break
		}
		start += len(indexes)
	}
}
