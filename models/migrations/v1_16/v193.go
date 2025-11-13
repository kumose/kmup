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
	"xorm.io/xorm"
)

func AddRepoIDForAttachment(x *xorm.Engine) error {
	type Attachment struct {
		ID         int64  `xorm:"pk autoincr"`
		UUID       string `xorm:"uuid UNIQUE"`
		RepoID     int64  `xorm:"INDEX"` // this should not be zero
		IssueID    int64  `xorm:"INDEX"` // maybe zero when creating
		ReleaseID  int64  `xorm:"INDEX"` // maybe zero when creating
		UploaderID int64  `xorm:"INDEX DEFAULT 0"`
	}
	if err := x.Sync(new(Attachment)); err != nil {
		return err
	}

	if _, err := x.Exec("UPDATE `attachment` set repo_id = (SELECT repo_id FROM `issue` WHERE `issue`.id = `attachment`.issue_id) WHERE `attachment`.issue_id > 0"); err != nil {
		return err
	}

	if _, err := x.Exec("UPDATE `attachment` set repo_id = (SELECT repo_id FROM `release` WHERE `release`.id = `attachment`.release_id) WHERE `attachment`.release_id > 0"); err != nil {
		return err
	}

	return nil
}
