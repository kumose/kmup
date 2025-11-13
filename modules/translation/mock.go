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

package translation

import (
	"fmt"
	"html/template"
	"strings"
)

// MockLocale provides a mocked locale without any translations
type MockLocale struct {
	Lang, LangName string // these fields are used directly in templates: ctx.Locale.Lang
}

var _ Locale = (*MockLocale)(nil)

func (l MockLocale) Language() string {
	return "en"
}

func (l MockLocale) TrString(s string, args ...any) string {
	return sprintAny(s, args...)
}

func (l MockLocale) Tr(s string, args ...any) template.HTML {
	return template.HTML(sprintAny(s, args...))
}

func (l MockLocale) TrN(cnt any, key1, keyN string, args ...any) template.HTML {
	return template.HTML(sprintAny(key1, args...))
}

func (l MockLocale) PrettyNumber(v any) string {
	return fmt.Sprint(v)
}

func sprintAny(s string, args ...any) string {
	if len(args) == 0 {
		return s
	}
	return s + ":" + fmt.Sprintf(strings.Repeat(",%v", len(args))[1:], args...)
}
