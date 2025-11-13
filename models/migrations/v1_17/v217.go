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
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm"
)

func AlterHookTaskTextFieldsToLongText(x *xorm.Engine) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if setting.Database.Type.IsMySQL() {
		if _, err := sess.Exec("ALTER TABLE `hook_task` CHANGE `payload_content` `payload_content` LONGTEXT, CHANGE `request_content` `request_content` LONGTEXT, change `response_content` `response_content` LONGTEXT"); err != nil {
			return err
		}
	}
	return sess.Commit()
}
