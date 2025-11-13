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

func TestPackAndUnpackData(t *testing.T) {
	s := "string"
	i := int64(4)
	f := float32(4.1)

	var s2 string
	var i2 int64
	var f2 float32

	data, err := PackData(s, i, f)
	assert.NoError(t, err)

	assert.NoError(t, UnpackData(data, &s2, &i2, &f2))
	assert.NoError(t, UnpackData(data, &s2))
	assert.Error(t, UnpackData(data, &i2))
	assert.Error(t, UnpackData(data, &s2, &f2))
}
