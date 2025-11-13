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

package rubygems

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimalEncoder(t *testing.T) {
	cases := []struct {
		Value    any
		Expected []byte
		Error    error
	}{
		{
			Value:    nil,
			Expected: []byte{4, 8, 0x30},
		},
		{
			Value:    true,
			Expected: []byte{4, 8, 'T'},
		},
		{
			Value:    false,
			Expected: []byte{4, 8, 'F'},
		},
		{
			Value:    0,
			Expected: []byte{4, 8, 'i', 0},
		},
		{
			Value:    1,
			Expected: []byte{4, 8, 'i', 6},
		},
		{
			Value:    -1,
			Expected: []byte{4, 8, 'i', 0xfa},
		},
		{
			Value:    0x1fffffff,
			Expected: []byte{4, 8, 'i', 4, 0xff, 0xff, 0xff, 0x1f},
		},
		{
			Value: 0x41000000,
			Error: ErrInvalidIntRange,
		},
		{
			Value:    "test",
			Expected: []byte{4, 8, 'I', '"', 9, 't', 'e', 's', 't', 6, ':', 6, 'E', 'T'},
		},
		{
			Value:    []int{1, 2},
			Expected: []byte{4, 8, '[', 7, 'i', 6, 'i', 7},
		},
		{
			Value: &RubyUserMarshal{
				Name:  "Test",
				Value: 4,
			},
			Expected: []byte{4, 8, 'U', ':', 9, 'T', 'e', 's', 't', 'i', 9},
		},
		{
			Value: &RubyUserDef{
				Name:  "Test",
				Value: 4,
			},
			Expected: []byte{4, 8, 'u', ':', 9, 'T', 'e', 's', 't', 9, 4, 8, 'i', 9},
		},
		{
			Value: &RubyObject{
				Name: "Test",
				Member: map[string]any{
					"test": 4,
				},
			},
			Expected: []byte{4, 8, 'o', ':', 9, 'T', 'e', 's', 't', 6, ':', 9, 't', 'e', 's', 't', 'i', 9},
		},
		{
			Value: &struct {
				Name string
			}{
				"test",
			},
			Error: ErrUnsupportedType,
		},
	}

	for i, c := range cases {
		var b bytes.Buffer
		err := NewMarshalEncoder(&b).Encode(c.Value)
		assert.ErrorIs(t, err, c.Error)
		assert.Equal(t, c.Expected, b.Bytes(), "case %d", i)
	}
}
