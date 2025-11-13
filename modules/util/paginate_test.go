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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginateSlice(t *testing.T) {
	stringSlice := []string{"a", "b", "c", "d", "e"}
	result, ok := PaginateSlice(stringSlice, 1, 2).([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, result)

	result, ok = PaginateSlice(stringSlice, 100, 2).([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{}, result)

	result, ok = PaginateSlice(stringSlice, 3, 2).([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"e"}, result)

	result, ok = PaginateSlice(stringSlice, 1, 0).([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, result)

	result, ok = PaginateSlice(stringSlice, 1, -1).([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, result)

	type Test struct {
		Val int
	}

	testVar := []*Test{{Val: 2}, {Val: 3}, {Val: 4}}
	testVar, ok = PaginateSlice(testVar, 1, 50).([]*Test)
	assert.True(t, ok)
	assert.Equal(t, []*Test{{Val: 2}, {Val: 3}, {Val: 4}}, testVar)

	testVar, ok = PaginateSlice(testVar, 2, 2).([]*Test)
	assert.True(t, ok)
	assert.Equal(t, []*Test{{Val: 4}}, testVar)
}
