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

package releasereopen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testReleaseReopener struct {
	count int
}

func (t *testReleaseReopener) ReleaseReopen() error {
	t.count++
	return nil
}

func TestManager(t *testing.T) {
	m := NewManager()

	t1 := &testReleaseReopener{}
	t2 := &testReleaseReopener{}
	t3 := &testReleaseReopener{}

	_ = m.Register(t1)
	c2 := m.Register(t2)
	_ = m.Register(t3)

	assert.NoError(t, m.ReleaseReopen())
	assert.Equal(t, 1, t1.count)
	assert.Equal(t, 1, t2.count)
	assert.Equal(t, 1, t3.count)

	c2()

	assert.NoError(t, m.ReleaseReopen())
	assert.Equal(t, 2, t1.count)
	assert.Equal(t, 1, t2.count)
	assert.Equal(t, 2, t3.count)
}
