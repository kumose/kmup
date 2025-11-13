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

package db

import (
	"strings"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/builder"
)

// BuildCaseInsensitiveLike returns a condition to check if the given value is like the given key case-insensitively.
// Handles especially SQLite correctly as UPPER there only transforms ASCII letters.
func BuildCaseInsensitiveLike(key, value string) builder.Cond {
	if setting.Database.Type.IsSQLite3() {
		return builder.Like{"UPPER(" + key + ")", util.ToUpperASCII(value)}
	}
	return builder.Like{"UPPER(" + key + ")", strings.ToUpper(value)}
}

// BuildCaseInsensitiveIn returns a condition to check if the given value is in the given values case-insensitively.
// Handles especially SQLite correctly as UPPER there only transforms ASCII letters.
func BuildCaseInsensitiveIn(key string, values []string) builder.Cond {
	uppers := make([]string, 0, len(values))
	if setting.Database.Type.IsSQLite3() {
		for _, value := range values {
			uppers = append(uppers, util.ToUpperASCII(value))
		}
	} else {
		for _, value := range values {
			uppers = append(uppers, strings.ToUpper(value))
		}
	}

	return builder.In("UPPER("+key+")", uppers)
}

// BuilderDialect returns the xorm.Builder dialect of the engine
func BuilderDialect() string {
	switch {
	case setting.Database.Type.IsMySQL():
		return builder.MYSQL
	case setting.Database.Type.IsSQLite3():
		return builder.SQLITE
	case setting.Database.Type.IsPostgreSQL():
		return builder.POSTGRES
	case setting.Database.Type.IsMSSQL():
		return builder.MSSQL
	default:
		return ""
	}
}
