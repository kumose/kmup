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

package v1_15

import (
	"github.com/kumose/kmup/models/migrations/base"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

func ConvertAvatarURLToText(x *xorm.Engine) error {
	dbType := x.Dialect().URI().DBType
	if dbType == schemas.SQLITE { // For SQLITE, varchar or char will always be represented as TEXT
		return nil
	}

	// Some oauth2 providers may give very long avatar urls (i.e. Google)
	return base.ModifyColumn(x, "external_login_user", &schemas.Column{
		Name: "avatar_url",
		SQLType: schemas.SQLType{
			Name: schemas.Text,
		},
		Nullable:       true,
		DefaultIsEmpty: true,
	})
}
