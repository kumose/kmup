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

package v1_24

import (
	"xorm.io/xorm"
)

func UpdateOwnerIDOfRepoLevelActionsTables(x *xorm.Engine) error {
	if _, err := x.Exec("UPDATE `action_runner` SET `owner_id` = 0 WHERE `repo_id` > 0 AND `owner_id` > 0"); err != nil {
		return err
	}
	if _, err := x.Exec("UPDATE `action_variable` SET `owner_id` = 0 WHERE `repo_id` > 0 AND `owner_id` > 0"); err != nil {
		return err
	}
	_, err := x.Exec("UPDATE `secret` SET `owner_id` = 0 WHERE `repo_id` > 0 AND `owner_id` > 0")
	return err
}
