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
	"fmt"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

// AddUniqueIndexForProjectIssue adds unique indexes for project issue table
func AddUniqueIndexForProjectIssue(x *xorm.Engine) error {
	// remove possible duplicated records in table project_issue
	type result struct {
		IssueID   int64
		ProjectID int64
		Cnt       int
	}
	var results []result
	if err := x.Select("issue_id, project_id, count(*) as cnt").
		Table("project_issue").
		GroupBy("issue_id, project_id").
		Having("count(*) > 1").
		Find(&results); err != nil {
		return err
	}
	for _, r := range results {
		if x.Dialect().URI().DBType == schemas.MSSQL {
			if _, err := x.Exec(fmt.Sprintf("delete from project_issue where id in (SELECT top %d id FROM project_issue WHERE issue_id = ? and project_id = ?)", r.Cnt-1), r.IssueID, r.ProjectID); err != nil {
				return err
			}
		} else {
			var ids []int64
			if err := x.SQL("SELECT id FROM project_issue WHERE issue_id = ? and project_id = ? limit ?", r.IssueID, r.ProjectID, r.Cnt-1).Find(&ids); err != nil {
				return err
			}
			if _, err := x.Table("project_issue").In("id", ids).Delete(); err != nil {
				return err
			}
		}
	}

	// add unique index for project_issue table
	type ProjectIssue struct { //revive:disable-line:exported
		ID        int64 `xorm:"pk autoincr"`
		IssueID   int64 `xorm:"INDEX unique(s)"`
		ProjectID int64 `xorm:"INDEX unique(s)"`
	}

	return x.Sync(new(ProjectIssue))
}
