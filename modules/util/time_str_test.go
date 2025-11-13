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

func TestTimeStr(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		// Test TimeEstimateParse
		tests := []struct {
			input  string
			output int64
			err    bool
		}{
			{"1h", 3600, false},
			{"1m", 60, false},
			{"1s", 1, false},
			{"1h 1m 1s", 3600 + 60 + 1, false},
			{"1d1x", 0, true},
		}
		for _, test := range tests {
			t.Run(test.input, func(t *testing.T) {
				output, err := TimeEstimateParse(test.input)
				if test.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, test.output, output)
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		tests := []struct {
			input  int64
			output string
		}{
			{3600, "1h"},
			{60, "1m"},
			{1, "1s"},
			{3600 + 1, "1h 1s"},
		}
		for _, test := range tests {
			t.Run(test.output, func(t *testing.T) {
				output := TimeEstimateString(test.input)
				assert.Equal(t, test.output, output)
			})
		}
	})
}
