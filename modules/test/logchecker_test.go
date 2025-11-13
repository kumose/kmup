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

package test

import (
	"testing"
	"time"

	"github.com/kumose/kmup/modules/log"

	"github.com/stretchr/testify/assert"
)

func TestLogChecker(t *testing.T) {
	lc, cleanup := NewLogChecker(log.DEFAULT)
	defer cleanup()

	lc.Filter("First", "Third").StopMark("End")
	log.Info("test")

	filtered, stopped := lc.Check(100 * time.Millisecond)
	assert.ElementsMatch(t, []bool{false, false}, filtered)
	assert.False(t, stopped)

	log.Info("First")
	filtered, stopped = lc.Check(100 * time.Millisecond)
	assert.ElementsMatch(t, []bool{true, false}, filtered)
	assert.False(t, stopped)

	log.Info("Second")
	filtered, stopped = lc.Check(100 * time.Millisecond)
	assert.ElementsMatch(t, []bool{true, false}, filtered)
	assert.False(t, stopped)

	log.Info("Third")
	filtered, stopped = lc.Check(100 * time.Millisecond)
	assert.ElementsMatch(t, []bool{true, true}, filtered)
	assert.False(t, stopped)

	log.Info("End")
	filtered, stopped = lc.Check(100 * time.Millisecond)
	assert.ElementsMatch(t, []bool{true, true}, filtered)
	assert.True(t, stopped)
}
