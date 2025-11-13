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
	"reflect"
)

type SliceUtils struct{}

func NewSliceUtils() *SliceUtils {
	return &SliceUtils{}
}

func (su *SliceUtils) Contains(s, v any) bool {
	if s == nil {
		return false
	}
	sv := reflect.ValueOf(s)
	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Array {
		panic(fmt.Sprintf("invalid type, expected slice or array, but got: %T", s))
	}
	for i := 0; i < sv.Len(); i++ {
		it := sv.Index(i)
		if !it.CanInterface() {
			continue
		}
		if it.Interface() == v {
			return true
		}
	}
	return false
}
