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
	"xorm.io/xorm"
)

func AddIssueResourceIndexTable(x *xorm.Engine) error {
	type ResourceIndex struct {
		GroupID  int64 `xorm:"pk"`
		MaxIndex int64 `xorm:"index"`
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if err := sess.Table("issue_index").Sync(new(ResourceIndex)); err != nil {
		return err
	}

	// Remove data we're goint to rebuild
	if _, err := sess.Table("issue_index").Where("1=1").Delete(&ResourceIndex{}); err != nil {
		return err
	}

	// Create current data for all repositories with issues and PRs
	if _, err := sess.Exec("INSERT INTO issue_index (group_id, max_index) " +
		"SELECT max_data.repo_id, max_data.max_index " +
		"FROM ( SELECT issue.repo_id AS repo_id, max(issue.`index`) AS max_index " +
		"FROM issue GROUP BY issue.repo_id) AS max_data"); err != nil {
		return err
	}

	return sess.Commit()
}
