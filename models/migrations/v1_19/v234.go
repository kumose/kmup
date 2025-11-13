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

package v1_19

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func CreatePackageCleanupRuleTable(x *xorm.Engine) error {
	type PackageCleanupRule struct {
		ID            int64              `xorm:"pk autoincr"`
		Enabled       bool               `xorm:"INDEX NOT NULL DEFAULT false"`
		OwnerID       int64              `xorm:"UNIQUE(s) INDEX NOT NULL DEFAULT 0"`
		Type          string             `xorm:"UNIQUE(s) INDEX NOT NULL"`
		KeepCount     int                `xorm:"NOT NULL DEFAULT 0"`
		KeepPattern   string             `xorm:"NOT NULL DEFAULT ''"`
		RemoveDays    int                `xorm:"NOT NULL DEFAULT 0"`
		RemovePattern string             `xorm:"NOT NULL DEFAULT ''"`
		MatchFullName bool               `xorm:"NOT NULL DEFAULT false"`
		CreatedUnix   timeutil.TimeStamp `xorm:"created NOT NULL DEFAULT 0"`
		UpdatedUnix   timeutil.TimeStamp `xorm:"updated NOT NULL DEFAULT 0"`
	}

	return x.Sync(new(PackageCleanupRule))
}
