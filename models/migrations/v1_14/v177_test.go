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

package v1_14

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/stretchr/testify/assert"
)

func Test_DeleteOrphanedIssueLabels(t *testing.T) {
	// Create the models used in the migration
	type IssueLabel struct {
		ID      int64 `xorm:"pk autoincr"`
		IssueID int64 `xorm:"UNIQUE(s)"`
		LabelID int64 `xorm:"UNIQUE(s)"`
	}

	type Label struct {
		ID              int64 `xorm:"pk autoincr"`
		RepoID          int64 `xorm:"INDEX"`
		OrgID           int64 `xorm:"INDEX"`
		Name            string
		Description     string
		Color           string `xorm:"VARCHAR(7)"`
		NumIssues       int
		NumClosedIssues int
		CreatedUnix     timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix     timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(IssueLabel), new(Label))
	if x == nil || t.Failed() {
		defer deferable()
		return
	}
	defer deferable()

	var issueLabels []*IssueLabel
	preMigration := map[int64]*IssueLabel{}
	postMigration := map[int64]*IssueLabel{}

	// Load issue labels that exist in the database pre-migration
	if err := x.Find(&issueLabels); err != nil {
		assert.NoError(t, err)
		return
	}
	for _, issueLabel := range issueLabels {
		preMigration[issueLabel.ID] = issueLabel
	}

	// Run the migration
	if err := DeleteOrphanedIssueLabels(x); err != nil {
		assert.NoError(t, err)
		return
	}

	// Load the remaining issue-labels
	issueLabels = issueLabels[:0]
	if err := x.Find(&issueLabels); err != nil {
		assert.NoError(t, err)
		return
	}
	for _, issueLabel := range issueLabels {
		postMigration[issueLabel.ID] = issueLabel
	}

	// Now test what is left
	if _, ok := postMigration[2]; ok {
		t.Errorf("Orphaned Label[2] survived the migration")
		return
	}

	if _, ok := postMigration[5]; ok {
		t.Errorf("Orphaned Label[5] survived the migration")
		return
	}

	for id, post := range postMigration {
		pre := preMigration[id]
		assert.Equal(t, pre, post, "migration changed issueLabel %d", id)
	}
}
