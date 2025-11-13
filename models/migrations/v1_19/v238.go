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

// AddUpdatedUnixToLFSMetaObject adds an updated column to the LFSMetaObject to allow for garbage collection
func AddUpdatedUnixToLFSMetaObject(x *xorm.Engine) error {
	// Drop the table introduced in `v211`, it's considered badly designed and doesn't look like to be used.
	// LFSMetaObject stores metadata for LFS tracked files.
	type LFSMetaObject struct {
		ID           int64              `xorm:"pk autoincr"`
		Oid          string             `json:"oid" xorm:"UNIQUE(s) INDEX NOT NULL"`
		Size         int64              `json:"size" xorm:"NOT NULL"`
		RepositoryID int64              `xorm:"UNIQUE(s) INDEX NOT NULL"`
		CreatedUnix  timeutil.TimeStamp `xorm:"created"`
		UpdatedUnix  timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	return x.Sync(new(LFSMetaObject))
}
