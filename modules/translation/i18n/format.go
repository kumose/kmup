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

package i18n

import (
	"errors"
	"fmt"
	"reflect"
)

// Format formats provided arguments for a given translated message
func Format(format string, args ...any) (msg string, err error) {
	if len(args) == 0 {
		return format, nil
	}

	fmtArgs := make([]any, 0, len(args))
	for _, arg := range args {
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.Slice {
			// Previously, we would accept Tr(lang, key, a, [b, c], d, [e, f]) as Sprintf(msg, a, b, c, d, e, f)
			// but this is an unstable behavior.
			//
			// So we restrict the accepted arguments to either:
			//
			// 1. Tr(lang, key, [slice-items]) as Sprintf(msg, items...)
			// 2. Tr(lang, key, args...) as Sprintf(msg, args...)
			if len(args) == 1 {
				for i := 0; i < val.Len(); i++ {
					fmtArgs = append(fmtArgs, val.Index(i).Interface())
				}
			} else {
				err = errors.New("arguments to i18n should not contain uncertain slices")
				break
			}
		} else {
			fmtArgs = append(fmtArgs, arg)
		}
	}
	return fmt.Sprintf(format, fmtArgs...), err
}
