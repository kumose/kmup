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

package unittest

import (
	"fmt"
	"reflect"
)

func fieldByName(v reflect.Value, field string) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName(field)
	if !f.IsValid() {
		panic(fmt.Errorf("can not read %s for %v", field, v))
	}
	return f
}

type reflectionValue struct {
	v reflect.Value
}

func reflectionWrap(v any) *reflectionValue {
	return &reflectionValue{v: reflect.ValueOf(v)}
}

func (rv *reflectionValue) int(field string) int {
	return int(fieldByName(rv.v, field).Int())
}

func (rv *reflectionValue) str(field string) string {
	return fieldByName(rv.v, field).String()
}

func (rv *reflectionValue) bool(field string) bool {
	return fieldByName(rv.v, field).Bool()
}
