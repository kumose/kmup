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

package v1_10

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddTaskTable(x *xorm.Engine) error {
	// TaskType defines task type
	type TaskType int

	// TaskStatus defines task status
	type TaskStatus int

	type Task struct {
		ID             int64
		DoerID         int64 `xorm:"index"` // operator
		OwnerID        int64 `xorm:"index"` // repo owner id, when creating, the repoID maybe zero
		RepoID         int64 `xorm:"index"`
		Type           TaskType
		Status         TaskStatus `xorm:"index"`
		StartTime      timeutil.TimeStamp
		EndTime        timeutil.TimeStamp
		PayloadContent string             `xorm:"TEXT"`
		Errors         string             `xorm:"TEXT"` // if task failed, saved the error reason
		Created        timeutil.TimeStamp `xorm:"created"`
	}

	type Repository struct {
		Status int `xorm:"NOT NULL DEFAULT 0"`
	}

	return x.Sync(new(Task), new(Repository))
}
