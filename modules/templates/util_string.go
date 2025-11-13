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

package templates

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/kumose/kmup/modules/util"
)

type StringUtils struct{}

var stringUtils = StringUtils{}

func NewStringUtils() *StringUtils {
	return &stringUtils
}

func (su *StringUtils) ToString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case template.HTML:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}

func (su *StringUtils) HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func (su *StringUtils) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (su *StringUtils) Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func (su *StringUtils) Join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func (su *StringUtils) Cut(s, sep string) []any {
	before, after, found := strings.Cut(s, sep)
	return []any{before, after, found}
}

func (su *StringUtils) EllipsisString(s string, maxLength int) string {
	return util.EllipsisDisplayString(s, maxLength)
}

func (su *StringUtils) ToUpper(s string) string {
	return strings.ToUpper(s)
}

func (su *StringUtils) TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}
