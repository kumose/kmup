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

package util

import "reflect"

// PaginateSlice cut a slice as per pagination options
// if page = 0 it do not paginate
func PaginateSlice(list any, page, pageSize int) any {
	if page <= 0 || pageSize <= 0 {
		return list
	}
	if reflect.TypeOf(list).Kind() != reflect.Slice {
		return list
	}

	listValue := reflect.ValueOf(list)

	page--

	if page*pageSize >= listValue.Len() {
		return listValue.Slice(listValue.Len(), listValue.Len()).Interface()
	}

	listValue = listValue.Slice(page*pageSize, listValue.Len())

	if listValue.Len() > pageSize {
		return listValue.Slice(0, pageSize).Interface()
	}

	return listValue.Interface()
}
