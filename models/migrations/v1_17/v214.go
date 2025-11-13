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

package v1_17

import (
	"xorm.io/xorm"
)

func AddAutoMergeTable(x *xorm.Engine) error {
	type MergeStyle string
	type PullAutoMerge struct {
		ID          int64      `xorm:"pk autoincr"`
		PullID      int64      `xorm:"UNIQUE"`
		DoerID      int64      `xorm:"NOT NULL"`
		MergeStyle  MergeStyle `xorm:"varchar(30)"`
		Message     string     `xorm:"LONGTEXT"`
		CreatedUnix int64      `xorm:"created"`
	}

	return x.Sync(&PullAutoMerge{})
}
