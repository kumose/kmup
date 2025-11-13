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
	"context"
	"fmt"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm"
)

func RenameTaskErrorsToMessage(x *xorm.Engine) error {
	type Task struct {
		Errors string `xorm:"TEXT"` // if task failed, saved the error reason
		Type   int
		Status int `xorm:"index"`
	}

	// This migration maybe rerun so that we should check if it has been run
	messageExist, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "task", "message")
	if err != nil {
		return err
	}

	if messageExist {
		errorsExist, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "task", "errors")
		if err != nil {
			return err
		}
		if !errorsExist {
			return nil
		}
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := sess.Sync(new(Task)); err != nil {
		return fmt.Errorf("error on Sync: %w", err)
	}

	if messageExist {
		// if both errors and message exist, drop message at first
		if err := base.DropTableColumns(sess, "task", "message"); err != nil {
			return err
		}
	}

	switch {
	case setting.Database.Type.IsMySQL():
		if _, err := sess.Exec("ALTER TABLE `task` CHANGE errors message text"); err != nil {
			return err
		}
	case setting.Database.Type.IsMSSQL():
		if _, err := sess.Exec("sp_rename 'task.errors', 'message', 'COLUMN'"); err != nil {
			return err
		}
	default:
		if _, err := sess.Exec("ALTER TABLE `task` RENAME COLUMN errors TO message"); err != nil {
			return err
		}
	}
	return sess.Commit()
}
