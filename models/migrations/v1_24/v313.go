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

package v1_24

import (
	"github.com/kumose/kmup/models/migrations/base"

	"xorm.io/xorm"
)

func MovePinOrderToTableIssuePin(x *xorm.Engine) error {
	type IssuePin struct {
		ID       int64 `xorm:"pk autoincr"`
		RepoID   int64 `xorm:"UNIQUE(s) NOT NULL"`
		IssueID  int64 `xorm:"UNIQUE(s) NOT NULL"`
		IsPull   bool  `xorm:"NOT NULL"`
		PinOrder int   `xorm:"DEFAULT 0"`
	}

	if err := x.Sync(new(IssuePin)); err != nil {
		return err
	}

	if _, err := x.Exec("INSERT INTO issue_pin (repo_id, issue_id, is_pull, pin_order) SELECT repo_id, id, is_pull, pin_order FROM issue WHERE pin_order > 0"); err != nil {
		return err
	}
	sess := x.NewSession()
	defer sess.Close()
	return base.DropTableColumns(sess, "issue", "pin_order")
}
