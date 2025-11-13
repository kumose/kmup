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

package v1_18

import (
	"xorm.io/xorm"
)

func CreateUserBadgesTable(x *xorm.Engine) error {
	type Badge struct {
		ID          int64 `xorm:"pk autoincr"`
		Description string
		ImageURL    string
	}

	type userBadge struct {
		ID      int64 `xorm:"pk autoincr"`
		BadgeID int64
		UserID  int64 `xorm:"INDEX"`
	}

	if err := x.Sync(new(Badge)); err != nil {
		return err
	}
	return x.Sync(new(userBadge))
}
