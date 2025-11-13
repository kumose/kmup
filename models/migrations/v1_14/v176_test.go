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

	"github.com/stretchr/testify/assert"
)

func Test_RemoveInvalidLabels(t *testing.T) {
	// Models used by the migration
	type Comment struct {
		ID           int64 `xorm:"pk autoincr"`
		Type         int   `xorm:"INDEX"`
		IssueID      int64 `xorm:"INDEX"`
		LabelID      int64
		ShouldRemain bool // <- Flag for testing the migration
	}

	type Issue struct {
		ID     int64 `xorm:"pk autoincr"`
		RepoID int64 `xorm:"INDEX UNIQUE(repo_index)"`
		Index  int64 `xorm:"UNIQUE(repo_index)"` // Index in one repository.
	}

	type Repository struct {
		ID        int64  `xorm:"pk autoincr"`
		OwnerID   int64  `xorm:"UNIQUE(s) index"`
		LowerName string `xorm:"UNIQUE(s) INDEX NOT NULL"`
	}

	type Label struct {
		ID     int64 `xorm:"pk autoincr"`
		RepoID int64 `xorm:"INDEX"`
		OrgID  int64 `xorm:"INDEX"`
	}

	type IssueLabel struct {
		ID           int64 `xorm:"pk autoincr"`
		IssueID      int64 `xorm:"UNIQUE(s)"`
		LabelID      int64 `xorm:"UNIQUE(s)"`
		ShouldRemain bool  // <- Flag for testing the migration
	}

	// load and prepare the test database
	x, deferable := base.PrepareTestEnv(t, 0, new(Comment), new(Issue), new(Repository), new(IssueLabel), new(Label))
	if x == nil || t.Failed() {
		defer deferable()
		return
	}
	defer deferable()

	var issueLabels []*IssueLabel
	ilPreMigration := map[int64]*IssueLabel{}
	ilPostMigration := map[int64]*IssueLabel{}

	var comments []*Comment
	comPreMigration := map[int64]*Comment{}
	comPostMigration := map[int64]*Comment{}

	// Get pre migration values
	if err := x.Find(&issueLabels); err != nil {
		t.Errorf("Unable to find issueLabels: %v", err)
		return
	}
	for _, issueLabel := range issueLabels {
		ilPreMigration[issueLabel.ID] = issueLabel
	}
	if err := x.Find(&comments); err != nil {
		t.Errorf("Unable to find comments: %v", err)
		return
	}
	for _, comment := range comments {
		comPreMigration[comment.ID] = comment
	}

	// Run the migration
	if err := RemoveInvalidLabels(x); err != nil {
		t.Errorf("unable to RemoveInvalidLabels: %v", err)
	}

	// Get the post migration values
	issueLabels = issueLabels[:0]
	if err := x.Find(&issueLabels); err != nil {
		t.Errorf("Unable to find issueLabels: %v", err)
		return
	}
	for _, issueLabel := range issueLabels {
		ilPostMigration[issueLabel.ID] = issueLabel
	}
	comments = comments[:0]
	if err := x.Find(&comments); err != nil {
		t.Errorf("Unable to find comments: %v", err)
		return
	}
	for _, comment := range comments {
		comPostMigration[comment.ID] = comment
	}

	// Finally test results of the migration
	for id, comment := range comPreMigration {
		post, ok := comPostMigration[id]
		if ok {
			if !comment.ShouldRemain {
				t.Errorf("Comment[%d] remained but should have been deleted", id)
			}
			assert.Equal(t, comment, post)
		} else if comment.ShouldRemain {
			t.Errorf("Comment[%d] was deleted but should have remained", id)
		}
	}

	for id, il := range ilPreMigration {
		post, ok := ilPostMigration[id]
		if ok {
			if !il.ShouldRemain {
				t.Errorf("IssueLabel[%d] remained but should have been deleted", id)
			}
			assert.Equal(t, il, post)
		} else if il.ShouldRemain {
			t.Errorf("IssueLabel[%d] was deleted but should have remained", id)
		}
	}
}
