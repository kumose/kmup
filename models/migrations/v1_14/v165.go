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

package v1_14

import (
	"github.com/kumose/kmup/models/migrations/base"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

func ConvertHookTaskTypeToVarcharAndTrim(x *xorm.Engine) error {
	dbType := x.Dialect().URI().DBType
	if dbType == schemas.SQLITE { // For SQLITE, varchar or char will always be represented as TEXT
		return nil
	}

	// HookTask: Typ string `xorm:"VARCHAR(16) index"`
	if err := base.ModifyColumn(x, "hook_task", &schemas.Column{
		Name: "typ",
		SQLType: schemas.SQLType{
			Name: "VARCHAR",
		},
		Length:         16,
		Nullable:       true, // To keep compatible as nullable
		DefaultIsEmpty: true,
	}); err != nil {
		return err
	}

	var hookTaskTrimSQL string
	if dbType == schemas.MSSQL {
		hookTaskTrimSQL = "UPDATE hook_task SET typ = RTRIM(LTRIM(typ))"
	} else {
		hookTaskTrimSQL = "UPDATE hook_task SET typ = TRIM(typ)"
	}
	if _, err := x.Exec(hookTaskTrimSQL); err != nil {
		return err
	}

	// Webhook: Type string `xorm:"VARCHAR(16) index"`
	if err := base.ModifyColumn(x, "webhook", &schemas.Column{
		Name: "type",
		SQLType: schemas.SQLType{
			Name: "VARCHAR",
		},
		Length:         16,
		Nullable:       true, // To keep compatible as nullable
		DefaultIsEmpty: true,
	}); err != nil {
		return err
	}

	var webhookTrimSQL string
	if dbType == schemas.MSSQL {
		webhookTrimSQL = "UPDATE webhook SET type = RTRIM(LTRIM(type))"
	} else {
		webhookTrimSQL = "UPDATE webhook SET type = TRIM(type)"
	}
	_, err := x.Exec(webhookTrimSQL)
	return err
}
