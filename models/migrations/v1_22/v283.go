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

func AddCombinedIndexToIssueUser(x *xorm.Engine) error {
	type OldIssueUser struct {
		IssueID int64
		UID     int64
		Cnt     int64
	}

	var duplicatedIssueUsers []OldIssueUser
	if err := x.SQL("select * from (select issue_id, uid, count(1) as cnt from issue_user group by issue_id, uid) a where a.cnt > 1").
		Find(&duplicatedIssueUsers); err != nil {
		return err
	}
	for _, issueUser := range duplicatedIssueUsers {
		if x.Dialect().URI().DBType == schemas.MSSQL {
			if _, err := x.Exec(fmt.Sprintf("delete from issue_user where id in (SELECT top %d id FROM issue_user WHERE issue_id = ? and uid = ?)", issueUser.Cnt-1), issueUser.IssueID, issueUser.UID); err != nil {
				return err
			}
		} else {
			var ids []int64
			if err := x.SQL("SELECT id FROM issue_user WHERE issue_id = ? and uid = ? limit ?", issueUser.IssueID, issueUser.UID, issueUser.Cnt-1).Find(&ids); err != nil {
				return err
			}
			if _, err := x.Table("issue_user").In("id", ids).Delete(); err != nil {
				return err
			}
		}
	}

	type IssueUser struct {
		UID     int64 `xorm:"INDEX unique(uid_to_issue)"` // User ID.
		IssueID int64 `xorm:"INDEX unique(uid_to_issue)"`
	}

	return x.Sync(&IssueUser{})
}
