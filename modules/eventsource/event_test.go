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

package eventsource

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_wrapNewlines(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		value  string
		output string
	}{
		{
			"check no new lines",
			"prefix: ",
			"value",
			"prefix: value\n",
		},
		{
			"check simple newline",
			"prefix: ",
			"value1\nvalue2",
			"prefix: value1\nprefix: value2\n",
		},
		{
			"check pathological newlines",
			"p: ",
			"\n1\n\n2\n3\n",
			"p: \np: 1\np: \np: 2\np: 3\np: \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			gotSum, err := wrapNewlines(w, []byte(tt.prefix), []byte(tt.value))
			require.NoError(t, err)

			assert.EqualValues(t, len(tt.output), gotSum)
			assert.Equal(t, tt.output, w.String())
		})
	}
}
