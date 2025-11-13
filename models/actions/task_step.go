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

package actions

import (
	"context"
	"time"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/timeutil"
)

// ActionTaskStep represents a step of ActionTask
type ActionTaskStep struct {
	ID        int64
	Name      string `xorm:"VARCHAR(255)"`
	TaskID    int64  `xorm:"index unique(task_index)"`
	Index     int64  `xorm:"index unique(task_index)"`
	RepoID    int64  `xorm:"index"`
	Status    Status `xorm:"index"`
	LogIndex  int64
	LogLength int64
	Started   timeutil.TimeStamp
	Stopped   timeutil.TimeStamp
	Created   timeutil.TimeStamp `xorm:"created"`
	Updated   timeutil.TimeStamp `xorm:"updated"`
}

func (step *ActionTaskStep) Duration() time.Duration {
	return calculateDuration(step.Started, step.Stopped, step.Status)
}

func init() {
	db.RegisterModel(new(ActionTaskStep))
}

func GetTaskStepsByTaskID(ctx context.Context, taskID int64) ([]*ActionTaskStep, error) {
	var steps []*ActionTaskStep
	return steps, db.GetEngine(ctx).Where("task_id=?", taskID).OrderBy("`index` ASC").Find(&steps)
}
