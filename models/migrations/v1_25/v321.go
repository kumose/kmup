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

package v1_25

import (
	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

func UseLongTextInSomeColumnsAndFixBugs(x *xorm.Engine) error {
	if !setting.Database.Type.IsMySQL() {
		return nil // Only mysql need to change from text to long text, for other databases, they are the same
	}

	if err := base.ModifyColumn(x, "review_state", &schemas.Column{
		Name: "updated_files",
		SQLType: schemas.SQLType{
			Name: "LONGTEXT",
		},
		Length:         0,
		Nullable:       false,
		DefaultIsEmpty: true,
	}); err != nil {
		return err
	}

	if err := base.ModifyColumn(x, "package_property", &schemas.Column{
		Name: "value",
		SQLType: schemas.SQLType{
			Name: "LONGTEXT",
		},
		Length:         0,
		Nullable:       false,
		DefaultIsEmpty: true,
	}); err != nil {
		return err
	}

	return base.ModifyColumn(x, "notice", &schemas.Column{
		Name: "description",
		SQLType: schemas.SQLType{
			Name: "LONGTEXT",
		},
		Length:         0,
		Nullable:       false,
		DefaultIsEmpty: true,
	})
}
