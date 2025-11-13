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

func TestSliceContainsString(t *testing.T) {
	assert.True(t, SliceContainsString([]string{"c", "b", "a", "b"}, "a"))
	assert.True(t, SliceContainsString([]string{"c", "b", "a", "b"}, "b"))
	assert.True(t, SliceContainsString([]string{"c", "b", "a", "b"}, "A", true))
	assert.True(t, SliceContainsString([]string{"C", "B", "A", "B"}, "a", true))

	assert.False(t, SliceContainsString([]string{"c", "b", "a", "b"}, "z"))
	assert.False(t, SliceContainsString([]string{"c", "b", "a", "b"}, "A"))
	assert.False(t, SliceContainsString([]string{}, "a"))
	assert.False(t, SliceContainsString(nil, "a"))
}

func TestSliceSortedEqual(t *testing.T) {
	assert.True(t, SliceSortedEqual([]int{2, 0, 2, 3}, []int{2, 0, 2, 3}))
	assert.True(t, SliceSortedEqual([]int{3, 0, 2, 2}, []int{2, 0, 2, 3}))
	assert.True(t, SliceSortedEqual([]int{}, []int{}))
	assert.True(t, SliceSortedEqual([]int(nil), nil))
	assert.True(t, SliceSortedEqual([]int(nil), []int{}))
	assert.True(t, SliceSortedEqual([]int{}, []int{}))

	assert.True(t, SliceSortedEqual([]string{"2", "0", "2", "3"}, []string{"2", "0", "2", "3"}))
	assert.True(t, SliceSortedEqual([]float64{2, 0, 2, 3}, []float64{2, 0, 2, 3}))
	assert.True(t, SliceSortedEqual([]bool{false, true, false}, []bool{false, true, false}))

	assert.False(t, SliceSortedEqual([]int{2, 0, 2}, []int{2, 0, 2, 3}))
	assert.False(t, SliceSortedEqual([]int{}, []int{2, 0, 2, 3}))
	assert.False(t, SliceSortedEqual(nil, []int{2, 0, 2, 3}))
	assert.False(t, SliceSortedEqual([]int{2, 0, 2, 4}, []int{2, 0, 2, 3}))
	assert.False(t, SliceSortedEqual([]int{2, 0, 0, 3}, []int{2, 0, 2, 3}))
}

func TestSliceRemoveAll(t *testing.T) {
	assert.ElementsMatch(t, []int{2, 2, 3}, SliceRemoveAll([]int{2, 0, 2, 3}, 0))
	assert.ElementsMatch(t, []int{0, 3}, SliceRemoveAll([]int{2, 0, 2, 3}, 2))
	assert.Empty(t, SliceRemoveAll([]int{0, 0, 0, 0}, 0))
	assert.ElementsMatch(t, []int{2, 0, 2, 3}, SliceRemoveAll([]int{2, 0, 2, 3}, 4))
	assert.Empty(t, SliceRemoveAll([]int{}, 0))
	assert.ElementsMatch(t, []int(nil), SliceRemoveAll([]int(nil), 0))
	assert.Empty(t, SliceRemoveAll([]int{}, 0))

	assert.ElementsMatch(t, []string{"2", "2", "3"}, SliceRemoveAll([]string{"2", "0", "2", "3"}, "0"))
	assert.ElementsMatch(t, []float64{2, 2, 3}, SliceRemoveAll([]float64{2, 0, 2, 3}, 0))
	assert.ElementsMatch(t, []bool{false, false}, SliceRemoveAll([]bool{false, true, false}, true))
}
