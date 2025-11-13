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

package v1_8

import (
	"fmt"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

func ChangeU2FCounterType(x *xorm.Engine) error {
	var err error

	switch x.Dialect().URI().DBType {
	case schemas.MYSQL:
		_, err = x.Exec("ALTER TABLE `u2f_registration` MODIFY `counter` BIGINT")
	case schemas.POSTGRES:
		_, err = x.Exec("ALTER TABLE `u2f_registration` ALTER COLUMN `counter` SET DATA TYPE bigint")
	case schemas.MSSQL:
		_, err = x.Exec("ALTER TABLE `u2f_registration` ALTER COLUMN `counter` BIGINT")
	}

	if err != nil {
		return fmt.Errorf("Error changing u2f_registration counter column type: %w", err)
	}

	return nil
}
