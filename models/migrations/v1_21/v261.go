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

func CreateVariableTable(x *xorm.Engine) error {
	type ActionVariable struct {
		ID          int64              `xorm:"pk autoincr"`
		OwnerID     int64              `xorm:"UNIQUE(owner_repo_name)"`
		RepoID      int64              `xorm:"INDEX UNIQUE(owner_repo_name)"`
		Name        string             `xorm:"UNIQUE(owner_repo_name) NOT NULL"`
		Data        string             `xorm:"LONGTEXT NOT NULL"`
		CreatedUnix timeutil.TimeStamp `xorm:"created NOT NULL"`
		UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
	}

	return x.Sync(new(ActionVariable))
}
