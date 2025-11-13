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

package system

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{FixtureFiles: []string{ /* load nothing */ }})
}

type testItem1 struct {
	Val1 string
	Val2 int
}

func (*testItem1) Name() string {
	return "test-item1"
}

type testItem2 struct {
	K string
}

func (*testItem2) Name() string {
	return "test-item2"
}

func TestAppStateDB(t *testing.T) {
	as := &DBStore{}

	item1 := new(testItem1)
	assert.NoError(t, as.Get(t.Context(), item1))
	assert.Empty(t, item1.Val1)
	assert.Equal(t, 0, item1.Val2)

	item1 = new(testItem1)
	item1.Val1 = "a"
	item1.Val2 = 2
	assert.NoError(t, as.Set(t.Context(), item1))

	item2 := new(testItem2)
	item2.K = "V"
	assert.NoError(t, as.Set(t.Context(), item2))

	item1 = new(testItem1)
	assert.NoError(t, as.Get(t.Context(), item1))
	assert.Equal(t, "a", item1.Val1)
	assert.Equal(t, 2, item1.Val2)

	item2 = new(testItem2)
	assert.NoError(t, as.Get(t.Context(), item2))
	assert.Equal(t, "V", item2.K)
}
