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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebounce(t *testing.T) {
	var c int64
	d := Debounce(50 * time.Millisecond)
	d(func() { atomic.AddInt64(&c, 1) })
	assert.EqualValues(t, 0, atomic.LoadInt64(&c))
	d(func() { atomic.AddInt64(&c, 1) })
	d(func() { atomic.AddInt64(&c, 1) })
	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, 1, atomic.LoadInt64(&c))
	d(func() { atomic.AddInt64(&c, 1) })
	assert.EqualValues(t, 1, atomic.LoadInt64(&c))
	d(func() { atomic.AddInt64(&c, 1) })
	d(func() { atomic.AddInt64(&c, 1) })
	d(func() { atomic.AddInt64(&c, 1) })
	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, 2, atomic.LoadInt64(&c))
}
