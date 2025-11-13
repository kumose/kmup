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

package v1_16

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

func Test_AddTableCommitStatusIndex(t *testing.T) {
	// Create the models used in the migration
	type CommitStatus struct {
		ID     int64  `xorm:"pk autoincr"`
		Index  int64  `xorm:"INDEX UNIQUE(repo_sha_index)"`
		RepoID int64  `xorm:"INDEX UNIQUE(repo_sha_index)"`
		SHA    string `xorm:"VARCHAR(64) NOT NULL INDEX UNIQUE(repo_sha_index)"`
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(CommitStatus))
	if x == nil || t.Failed() {
		defer deferable()
		return
	}
	defer deferable()

	// Run the migration
	if err := AddTableCommitStatusIndex(x); err != nil {
		assert.NoError(t, err)
		return
	}

	type CommitStatusIndex struct {
		ID       int64
		RepoID   int64  `xorm:"unique(repo_sha)"`
		SHA      string `xorm:"unique(repo_sha)"`
		MaxIndex int64  `xorm:"index"`
	}

	start := 0
	const batchSize = 1000
	for {
		indexes := make([]CommitStatusIndex, 0, batchSize)
		err := x.Table("commit_status_index").Limit(batchSize, start).Find(&indexes)
		assert.NoError(t, err)

		for _, idx := range indexes {
			var maxIndex int
			has, err := x.SQL("SELECT max(`index`) FROM commit_status WHERE repo_id = ? AND sha = ?", idx.RepoID, idx.SHA).Get(&maxIndex)
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
