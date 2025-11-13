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
	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func UpdateReactionConstraint(x *xorm.Engine) error {
	// Reaction represents a reactions on issues and comments.
	type Reaction struct {
		ID               int64              `xorm:"pk autoincr"`
		Type             string             `xorm:"INDEX UNIQUE(s) NOT NULL"`
		IssueID          int64              `xorm:"INDEX UNIQUE(s) NOT NULL"`
		CommentID        int64              `xorm:"INDEX UNIQUE(s)"`
		UserID           int64              `xorm:"INDEX UNIQUE(s) NOT NULL"`
		OriginalAuthorID int64              `xorm:"INDEX UNIQUE(s) NOT NULL DEFAULT(0)"`
		OriginalAuthor   string             `xorm:"INDEX UNIQUE(s)"`
		CreatedUnix      timeutil.TimeStamp `xorm:"INDEX created"`
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if err := base.RecreateTable(sess, &Reaction{}); err != nil {
		return err
	}

	return sess.Commit()
}
