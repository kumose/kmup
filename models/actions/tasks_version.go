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
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/timeutil"
)

// ActionTasksVersion
// If both ownerID and repoID is zero, its scope is global.
// If ownerID is not zero and repoID is zero, its scope is org (there is no user-level runner currently).
// If ownerID is zero and repoID is not zero, its scope is repo.
type ActionTasksVersion struct {
	ID          int64 `xorm:"pk autoincr"`
	OwnerID     int64 `xorm:"UNIQUE(owner_repo)"`
	RepoID      int64 `xorm:"INDEX UNIQUE(owner_repo)"`
	Version     int64
	CreatedUnix timeutil.TimeStamp `xorm:"created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
}

func init() {
	db.RegisterModel(new(ActionTasksVersion))
}

func GetTasksVersionByScope(ctx context.Context, ownerID, repoID int64) (int64, error) {
	var tasksVersion ActionTasksVersion
	has, err := db.GetEngine(ctx).Where("owner_id = ? AND repo_id = ?", ownerID, repoID).Get(&tasksVersion)
	if err != nil {
		return 0, err
	} else if !has {
		return 0, nil
	}
	return tasksVersion.Version, err
}

func insertTasksVersion(ctx context.Context, ownerID, repoID int64) (*ActionTasksVersion, error) {
	tasksVersion := &ActionTasksVersion{
		OwnerID: ownerID,
		RepoID:  repoID,
		Version: 1,
	}
	if _, err := db.GetEngine(ctx).Insert(tasksVersion); err != nil {
		return nil, err
	}
	return tasksVersion, nil
}

func increaseTasksVersionByScope(ctx context.Context, ownerID, repoID int64) error {
	result, err := db.GetEngine(ctx).Exec("UPDATE action_tasks_version SET version = version + 1 WHERE owner_id = ? AND repo_id = ?", ownerID, repoID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		// if update sql does not affect any rows, the database may be broken,
		// so re-insert the row of version data here.
		if _, err := insertTasksVersion(ctx, ownerID, repoID); err != nil {
			return err
		}
	}

	return nil
}

func IncreaseTaskVersion(ctx context.Context, ownerID, repoID int64) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		// 1. increase global
		if err := increaseTasksVersionByScope(ctx, 0, 0); err != nil {
			log.Error("IncreaseTasksVersionByScope(Global): %v", err)
			return err
		}

		// 2. increase owner
		if ownerID > 0 {
			if err := increaseTasksVersionByScope(ctx, ownerID, 0); err != nil {
				log.Error("IncreaseTasksVersionByScope(Owner): %v", err)
				return err
			}
		}

		// 3. increase repo
		if repoID > 0 {
			if err := increaseTasksVersionByScope(ctx, 0, repoID); err != nil {
				log.Error("IncreaseTasksVersionByScope(Repo): %v", err)
				return err
			}
		}

		return nil
	})
}
