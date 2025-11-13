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
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UseLongTextInSomeColumnsAndFixBugs(t *testing.T) {
	if !setting.Database.Type.IsMySQL() {
		t.Skip("Only MySQL needs to change from TEXT to LONGTEXT")
	}

	type ReviewState struct {
		ID           int64              `xorm:"pk autoincr"`
		UserID       int64              `xorm:"NOT NULL UNIQUE(pull_commit_user)"`
		PullID       int64              `xorm:"NOT NULL INDEX UNIQUE(pull_commit_user) DEFAULT 0"` // Which PR was the review on?
		CommitSHA    string             `xorm:"NOT NULL VARCHAR(64) UNIQUE(pull_commit_user)"`     // Which commit was the head commit for the review?
		UpdatedFiles map[string]int     `xorm:"NOT NULL TEXT JSON"`                                // Stores for each of the changed files of a PR whether they have been viewed, changed since last viewed, or not viewed
		UpdatedUnix  timeutil.TimeStamp `xorm:"updated"`                                           // Is an accurate indicator of the order of commits as we do not expect it to be possible to make reviews on previous commits
	}

	type PackageProperty struct {
		ID      int64  `xorm:"pk autoincr"`
		RefType int    `xorm:"INDEX NOT NULL"`
		RefID   int64  `xorm:"INDEX NOT NULL"`
		Name    string `xorm:"INDEX NOT NULL"`
		Value   string `xorm:"TEXT NOT NULL"`
	}

	type Notice struct {
		ID          int64 `xorm:"pk autoincr"`
		Type        int
		Description string             `xorm:"TEXT"`
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	}

	// Prepare and load the testing database
	x, deferrable := base.PrepareTestEnv(t, 0, new(ReviewState), new(PackageProperty), new(Notice))
	defer deferrable()

	require.NoError(t, UseLongTextInSomeColumnsAndFixBugs(x))

	tables := base.LoadTableSchemasMap(t, x)
	table := tables["review_state"]
	column := table.GetColumn("updated_files")
	assert.Equal(t, "LONGTEXT", column.SQLType.Name)

	table = tables["package_property"]
	column = table.GetColumn("value")
	assert.Equal(t, "LONGTEXT", column.SQLType.Name)

	table = tables["notice"]
	column = table.GetColumn("description")
	assert.Equal(t, "LONGTEXT", column.SQLType.Name)
}
