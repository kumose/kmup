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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterMapUnique(t *testing.T) {
	result := FilterSlice([]int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	}, func(i int) (int, bool) {
		switch i {
		case 0:
			return 0, true // included later
		case 1:
			return 0, true // duplicate of previous (should be ignored)
		case 2:
			return 2, false // not included
		default:
			return i, true
		}
	})
	assert.Equal(t, []int{0, 3, 4, 5, 6, 7, 8, 9}, result)
}
