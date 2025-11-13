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
	"xorm.io/xorm/schemas"
)

func ChangeReviewContentToText(x *xorm.Engine) error {
	switch x.Dialect().URI().DBType {
	case schemas.MYSQL:
		_, err := x.Exec("ALTER TABLE review MODIFY COLUMN content TEXT")
		return err
	case schemas.ORACLE:
		_, err := x.Exec("ALTER TABLE review MODIFY content TEXT")
		return err
	case schemas.MSSQL:
		_, err := x.Exec("ALTER TABLE review ALTER COLUMN content TEXT")
		return err
	case schemas.POSTGRES:
		_, err := x.Exec("ALTER TABLE review ALTER COLUMN content TYPE TEXT")
		return err
	default:
		// SQLite doesn't support ALTER COLUMN, and it seem to already make String to _TEXT_ default so no migration needed
		return nil
	}
}
