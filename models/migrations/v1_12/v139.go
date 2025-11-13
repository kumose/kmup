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

package v1_12

import (
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm"
)

func PrependRefsHeadsToIssueRefs(x *xorm.Engine) error {
	var query string

	switch {
	case setting.Database.Type.IsMSSQL():
		query = "UPDATE `issue` SET `ref` = 'refs/heads/' + `ref` WHERE `ref` IS NOT NULL AND `ref` <> '' AND `ref` NOT LIKE 'refs/%'"
	case setting.Database.Type.IsMySQL():
		query = "UPDATE `issue` SET `ref` = CONCAT('refs/heads/', `ref`) WHERE `ref` IS NOT NULL AND `ref` <> '' AND `ref` NOT LIKE 'refs/%';"
	default:
		query = "UPDATE `issue` SET `ref` = 'refs/heads/' || `ref` WHERE `ref` IS NOT NULL AND `ref` <> '' AND `ref` NOT LIKE 'refs/%'"
	}

	_, err := x.Exec(query)
	return err
}
