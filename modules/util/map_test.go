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

func TestGetMapValueOrDefault(t *testing.T) {
	testMap := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": nil,
	}

	assert.Equal(t, "value1", GetMapValueOrDefault(testMap, "key1", "default"))
	assert.Equal(t, 42, GetMapValueOrDefault(testMap, "key2", 0))

	assert.Equal(t, "default", GetMapValueOrDefault(testMap, "key4", "default"))
	assert.Equal(t, 100, GetMapValueOrDefault(testMap, "key5", 100))

	assert.Equal(t, "default", GetMapValueOrDefault(testMap, "key3", "default"))
}
