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

package v1_20

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

type Action struct {
	UserID      int64 // Receiver user id.
	ActUserID   int64 // Action user id.
	RepoID      int64
	IsDeleted   bool               `xorm:"NOT NULL DEFAULT false"`
	IsPrivate   bool               `xorm:"NOT NULL DEFAULT false"`
	CreatedUnix timeutil.TimeStamp `xorm:"created"`
}

// TableName sets the name of this table
func (a *Action) TableName() string {
	return "action"
}

// TableIndices implements xorm's TableIndices interface
func (a *Action) TableIndices() []*schemas.Index {
	repoIndex := schemas.NewIndex("r_u_d", schemas.IndexType)
	repoIndex.AddColumn("repo_id", "user_id", "is_deleted")

	actUserIndex := schemas.NewIndex("au_r_c_u_d", schemas.IndexType)
	actUserIndex.AddColumn("act_user_id", "repo_id", "created_unix", "user_id", "is_deleted")

	cudIndex := schemas.NewIndex("c_u_d", schemas.IndexType)
	cudIndex.AddColumn("created_unix", "user_id", "is_deleted")

	indices := []*schemas.Index{actUserIndex, repoIndex, cudIndex}

	return indices
}

func ImproveActionTableIndices(x *xorm.Engine) error {
	return x.Sync(new(Action))
}
