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

func Test_AddRepoIDForAttachment(t *testing.T) {
	type Attachment struct {
		ID         int64  `xorm:"pk autoincr"`
		UUID       string `xorm:"uuid UNIQUE"`
		IssueID    int64  `xorm:"INDEX"` // maybe zero when creating
		ReleaseID  int64  `xorm:"INDEX"` // maybe zero when creating
		UploaderID int64  `xorm:"INDEX DEFAULT 0"`
	}

	type Issue struct {
		ID     int64
		RepoID int64
	}

	type Release struct {
		ID     int64
		RepoID int64
	}

	// Prepare and load the testing database
	x, deferrable := base.PrepareTestEnv(t, 0, new(Attachment), new(Issue), new(Release))
	defer deferrable()
	if x == nil || t.Failed() {
		return
	}

	// Run the migration
	if err := AddRepoIDForAttachment(x); err != nil {
		assert.NoError(t, err)
		return
	}

	type NewAttachment struct {
		ID         int64  `xorm:"pk autoincr"`
		UUID       string `xorm:"uuid UNIQUE"`
		RepoID     int64  `xorm:"INDEX"` // this should not be zero
		IssueID    int64  `xorm:"INDEX"` // maybe zero when creating
		ReleaseID  int64  `xorm:"INDEX"` // maybe zero when creating
		UploaderID int64  `xorm:"INDEX DEFAULT 0"`
	}

	var issueAttachments []*NewAttachment
	err := x.Table("attachment").Where("issue_id > 0").Find(&issueAttachments)
	assert.NoError(t, err)
	for _, attach := range issueAttachments {
		assert.Positive(t, attach.RepoID)
		assert.Positive(t, attach.IssueID)
		var issue Issue
		has, err := x.ID(attach.IssueID).Get(&issue)
		assert.NoError(t, err)
		assert.True(t, has)
		assert.Equal(t, attach.RepoID, issue.RepoID)
	}

	var releaseAttachments []*NewAttachment
	err = x.Table("attachment").Where("release_id > 0").Find(&releaseAttachments)
	assert.NoError(t, err)
	for _, attach := range releaseAttachments {
		assert.Positive(t, attach.RepoID)
		assert.Positive(t, attach.ReleaseID)
		var release Release
		has, err := x.ID(attach.ReleaseID).Get(&release)
		assert.NoError(t, err)
		assert.True(t, has)
		assert.Equal(t, attach.RepoID, release.RepoID)
	}
}
