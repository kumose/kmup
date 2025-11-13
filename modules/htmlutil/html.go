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

package htmlutil

import (
	"fmt"
	"html/template"
	"slices"
	"strings"
)

// ParseSizeAndClass get size and class from string with default values
// If present, "others" expects the new size first and then the classes to use
func ParseSizeAndClass(defaultSize int, defaultClass string, others ...any) (int, string) {
	size := defaultSize
	if len(others) >= 1 {
		if v, ok := others[0].(int); ok && v != 0 {
			size = v
		}
	}
	class := defaultClass
	if len(others) >= 2 {
		if v, ok := others[1].(string); ok && v != "" {
			if class != "" {
				class += " "
			}
			class += v
		}
	}
	return size, class
}

func HTMLFormat(s template.HTML, rawArgs ...any) template.HTML {
	if !strings.Contains(string(s), "%") || len(rawArgs) == 0 {
		panic("HTMLFormat requires one or more arguments")
	}
	args := slices.Clone(rawArgs)
	for i, v := range args {
		switch v := v.(type) {
		case nil, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, template.HTML:
			// for most basic types (including template.HTML which is safe), just do nothing and use it
		case string:
			args[i] = template.HTMLEscapeString(v)
		case template.URL:
			args[i] = template.HTMLEscapeString(string(v))
		case fmt.Stringer:
			args[i] = template.HTMLEscapeString(v.String())
		default:
			args[i] = template.HTMLEscapeString(fmt.Sprint(v))
		}
	}
	return template.HTML(fmt.Sprintf(string(s), args...))
}
