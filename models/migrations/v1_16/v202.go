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

package v1_16

import (
	"fmt"

	"xorm.io/xorm"
)

func CreateUserSettingsTable(x *xorm.Engine) error {
	type UserSetting struct {
		ID           int64  `xorm:"pk autoincr"`
		UserID       int64  `xorm:"index unique(key_userid)"`              // to load all of someone's settings
		SettingKey   string `xorm:"varchar(255) index unique(key_userid)"` // ensure key is always lowercase
		SettingValue string `xorm:"text"`
	}
	if err := x.Sync(new(UserSetting)); err != nil {
		return fmt.Errorf("sync2: %w", err)
	}
	return nil
}
