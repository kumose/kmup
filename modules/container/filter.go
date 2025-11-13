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

package container

import "slices"

// FilterSlice ranges over the slice and calls include() for each element.
// If the second returned value is true, the first returned value will be included in the resulting
// slice (after deduplication).
func FilterSlice[E any, T comparable](s []E, include func(E) (T, bool)) []T {
	filtered := make([]T, 0, len(s)) // slice will be clipped before returning
	seen := make(map[T]bool, len(s))
	for i := range s {
		if v, ok := include(s[i]); ok && !seen[v] {
			filtered = append(filtered, v)
			seen[v] = true
		}
	}
	return slices.Clip(filtered)
}
