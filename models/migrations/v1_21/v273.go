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

package v1_21

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddActionScheduleTable(x *xorm.Engine) error {
	type ActionSchedule struct {
		ID            int64
		Title         string
		Specs         []string
		RepoID        int64 `xorm:"index"`
		OwnerID       int64 `xorm:"index"`
		WorkflowID    string
		TriggerUserID int64
		Ref           string
		CommitSHA     string
		Event         string
		EventPayload  string `xorm:"LONGTEXT"`
		Content       []byte
		Created       timeutil.TimeStamp `xorm:"created"`
		Updated       timeutil.TimeStamp `xorm:"updated"`
	}

	type ActionScheduleSpec struct {
		ID         int64
		RepoID     int64 `xorm:"index"`
		ScheduleID int64 `xorm:"index"`
		Spec       string
		Next       timeutil.TimeStamp `xorm:"index"`
		Prev       timeutil.TimeStamp

		Created timeutil.TimeStamp `xorm:"created"`
		Updated timeutil.TimeStamp `xorm:"updated"`
	}

	return x.Sync(
		new(ActionSchedule),
		new(ActionScheduleSpec),
	)
}
