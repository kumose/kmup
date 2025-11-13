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

package v1_11

import (
	"xorm.io/xorm"
)

func AddTeamIncludesAllRepositories(x *xorm.Engine) error {
	type Team struct {
		ID                      int64 `xorm:"pk autoincr"`
		IncludesAllRepositories bool  `xorm:"NOT NULL DEFAULT false"`
	}

	if err := x.Sync(new(Team)); err != nil {
		return err
	}

	_, err := x.Exec("UPDATE `team` SET `includes_all_repositories` = ? WHERE `name`=?",
		true, "Owners")
	return err
}
