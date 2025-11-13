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

func TestSet(t *testing.T) {
	s := make(Set[string])

	assert.True(t, s.Add("key1"))
	assert.False(t, s.Add("key1"))
	assert.True(t, s.Add("key2"))

	assert.True(t, s.Contains("key1"))
	assert.True(t, s.Contains("key2"))
	assert.True(t, s.Contains("key1", "key2"))
	assert.False(t, s.Contains("key3"))
	assert.False(t, s.Contains("key1", "key3"))

	assert.True(t, s.Remove("key2"))
	assert.False(t, s.Contains("key2"))

	assert.False(t, s.Remove("key3"))

	s.AddMultiple("key4", "key5")
	assert.True(t, s.Contains("key4"))
	assert.True(t, s.Contains("key5"))

	s = SetOf("key6", "key7")
	assert.False(t, s.Contains("key1"))
	assert.True(t, s.Contains("key6"))
	assert.True(t, s.Contains("key7"))
}
