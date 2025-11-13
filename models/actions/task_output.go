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

	"github.com/kumose/kmup/models/db"
)

// ActionTaskOutput represents an output of ActionTask.
// So the outputs are bound to a task, that means when a completed job has been rerun,
// the outputs of the job will be reset because the task is new.
// It's by design, to avoid the outputs of the old task to be mixed with the new task.
type ActionTaskOutput struct {
	ID          int64
	TaskID      int64  `xorm:"INDEX UNIQUE(task_id_output_key)"`
	OutputKey   string `xorm:"VARCHAR(255) UNIQUE(task_id_output_key)"`
	OutputValue string `xorm:"MEDIUMTEXT"`
}

func init() {
	db.RegisterModel(new(ActionTaskOutput))
}

// FindTaskOutputByTaskID returns the outputs of the task.
func FindTaskOutputByTaskID(ctx context.Context, taskID int64) ([]*ActionTaskOutput, error) {
	var outputs []*ActionTaskOutput
	return outputs, db.GetEngine(ctx).Where("task_id=?", taskID).Find(&outputs)
}

// FindTaskOutputKeyByTaskID returns the keys of the outputs of the task.
func FindTaskOutputKeyByTaskID(ctx context.Context, taskID int64) ([]string, error) {
	var keys []string
	return keys, db.GetEngine(ctx).Table(ActionTaskOutput{}).Where("task_id=?", taskID).Cols("output_key").Find(&keys)
}

// InsertTaskOutputIfNotExist inserts a new task output if it does not exist.
func InsertTaskOutputIfNotExist(ctx context.Context, taskID int64, key, value string) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		sess := db.GetEngine(ctx)
		if exist, err := sess.Exist(&ActionTaskOutput{TaskID: taskID, OutputKey: key}); err != nil {
			return err
		} else if exist {
			return nil
		}
		_, err := sess.Insert(&ActionTaskOutput{
			TaskID:      taskID,
			OutputKey:   key,
			OutputValue: value,
		})
		return err
	})
}
