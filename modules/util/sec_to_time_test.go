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

func TestSecToHours(t *testing.T) {
	second := int64(1)
	minute := 60 * second
	hour := 60 * minute
	day := 24 * hour

	assert.Equal(t, "1 minute", SecToHours(minute+6*second))
	assert.Equal(t, "1 hour", SecToHours(hour))
	assert.Equal(t, "1 hour", SecToHours(hour+second))
	assert.Equal(t, "14 hours 33 minutes", SecToHours(14*hour+33*minute+30*second))
	assert.Equal(t, "156 hours 30 minutes", SecToHours(6*day+12*hour+30*minute+18*second))
	assert.Equal(t, "98 hours 16 minutes", SecToHours(4*day+2*hour+16*minute+58*second))
	assert.Equal(t, "672 hours", SecToHours(4*7*day))
	assert.Equal(t, "1 second", SecToHours(1))
	assert.Equal(t, "2 seconds", SecToHours(2))
	assert.Empty(t, SecToHours(nil)) // old behavior, empty means no output
}
