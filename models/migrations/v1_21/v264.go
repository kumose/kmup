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
	"context"
	"errors"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddBranchTable(x *xorm.Engine) error {
	type Branch struct {
		ID            int64
		RepoID        int64  `xorm:"UNIQUE(s)"`
		Name          string `xorm:"UNIQUE(s) NOT NULL"`
		CommitID      string
		CommitMessage string `xorm:"TEXT"`
		PusherID      int64
		IsDeleted     bool `xorm:"index"`
		DeletedByID   int64
		DeletedUnix   timeutil.TimeStamp `xorm:"index"`
		CommitTime    timeutil.TimeStamp // The commit
		CreatedUnix   timeutil.TimeStamp `xorm:"created"`
		UpdatedUnix   timeutil.TimeStamp `xorm:"updated"`
	}

	if err := x.Sync(new(Branch)); err != nil {
		return err
	}

	if exist, err := x.IsTableExist("deleted_branches"); err != nil {
		return err
	} else if !exist {
		return nil
	}

	type DeletedBranch struct {
		ID          int64
		RepoID      int64  `xorm:"index UNIQUE(s)"`
		Name        string `xorm:"UNIQUE(s) NOT NULL"`
		Commit      string
		DeletedByID int64
		DeletedUnix timeutil.TimeStamp
	}

	var adminUserID int64
	has, err := x.Table("user").
		Select("id").
		Where("is_admin=?", true).
		Asc("id"). // Reliably get the admin with the lowest ID.
		Get(&adminUserID)
	if err != nil {
		return err
	} else if !has {
		return errors.New("no admin user found")
	}

	branches := make([]Branch, 0, 100)
	if err := db.Iterate(context.Background(), nil, func(ctx context.Context, deletedBranch *DeletedBranch) error {
		branches = append(branches, Branch{
			RepoID:      deletedBranch.RepoID,
			Name:        deletedBranch.Name,
			CommitID:    deletedBranch.Commit,
			PusherID:    adminUserID,
			IsDeleted:   true,
			DeletedByID: deletedBranch.DeletedByID,
			DeletedUnix: deletedBranch.DeletedUnix,
		})
		if len(branches) >= 100 {
			_, err := x.Insert(&branches)
			if err != nil {
				return err
			}
			branches = branches[:0]
		}
		return nil
	}); err != nil {
		return err
	}

	if len(branches) > 0 {
		if _, err := x.Insert(&branches); err != nil {
			return err
		}
	}

	return x.DropTables(new(DeletedBranch))
}
