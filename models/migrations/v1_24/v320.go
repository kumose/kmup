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
	"github.com/kumose/kmup/modules/json"

	"xorm.io/xorm"
)

func MigrateSkipTwoFactor(x *xorm.Engine) error {
	type LoginSource struct {
		TwoFactorPolicy string `xorm:"two_factor_policy NOT NULL DEFAULT ''"`
	}
	_, err := x.SyncWithOptions(
		xorm.SyncOptions{
			IgnoreConstrains: true,
			IgnoreIndices:    true,
		},
		new(LoginSource),
	)
	if err != nil {
		return err
	}

	type LoginSourceSimple struct {
		ID  int64
		Cfg string
	}

	var loginSources []LoginSourceSimple
	err = x.Table("login_source").Find(&loginSources)
	if err != nil {
		return err
	}

	for _, source := range loginSources {
		if source.Cfg == "" {
			continue
		}

		var cfg map[string]any
		err = json.Unmarshal([]byte(source.Cfg), &cfg)
		if err != nil {
			return err
		}

		if cfg["SkipLocalTwoFA"] == true {
			_, err = x.Exec("UPDATE login_source SET two_factor_policy = 'skip' WHERE id = ?", source.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
