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
	"fmt"

	"xorm.io/xorm"
)

// DeleteOrphanedIssueLabels looks through the database for issue_labels where the label no longer exists and deletes them.
func DeleteOrphanedIssueLabels(x *xorm.Engine) error {
	type IssueLabel struct {
		ID      int64 `xorm:"pk autoincr"`
		IssueID int64 `xorm:"UNIQUE(s)"`
		LabelID int64 `xorm:"UNIQUE(s)"`
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := sess.Sync(new(IssueLabel)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	if _, err := sess.Exec(`DELETE FROM issue_label WHERE issue_label.id IN (
		SELECT ill.id FROM (
			SELECT il.id
			FROM issue_label AS il
				LEFT JOIN label ON il.label_id = label.id
			WHERE
				label.id IS NULL
		) AS ill)`); err != nil {
		return err
	}

	return sess.Commit()
}
