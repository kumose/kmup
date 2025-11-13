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

package v1_7

import (
	"xorm.io/builder"
	"xorm.io/xorm"
)

func ClearNonusedData(x *xorm.Engine) error {
	condDelete := func(colName string) builder.Cond {
		return builder.NotIn(colName, builder.Select("id").From("`user`"))
	}

	if _, err := x.Exec(builder.Delete(condDelete("uid")).From("team_user")); err != nil {
		return err
	}

	if _, err := x.Exec(builder.Delete(condDelete("user_id")).From("collaboration")); err != nil {
		return err
	}

	if _, err := x.Exec(builder.Delete(condDelete("user_id")).From("stopwatch")); err != nil {
		return err
	}

	if _, err := x.Exec(builder.Delete(condDelete("owner_id")).From("gpg_key")); err != nil {
		return err
	}
	return nil
}
