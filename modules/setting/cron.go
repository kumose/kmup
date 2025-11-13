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

package setting

import "reflect"

// GetCronSettings maps the cron subsection to the provided config
func GetCronSettings(name string, config any) (any, error) {
	return getCronSettings(CfgProvider, name, config)
}

func getCronSettings(rootCfg ConfigProvider, name string, config any) (any, error) {
	if err := rootCfg.Section("cron." + name).MapTo(config); err != nil {
		return config, err
	}

	typ := reflect.TypeOf(config).Elem()
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		tpField := typ.Field(i)
		if tpField.Type.Kind() == reflect.Struct && tpField.Anonymous {
			if err := rootCfg.Section("cron." + name).MapTo(field.Addr().Interface()); err != nil {
				return config, err
			}
		}
	}

	return config, nil
}
