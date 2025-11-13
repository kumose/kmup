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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeFlashErrorString(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "no error",
			arg:  "",
			want: "",
		},
		{
			name: "normal error",
			arg:  "can not open file: \"abc.exe\"",
			want: "can not open file: &#34;abc.exe&#34;",
		},
		{
			name: "line break error",
			arg:  "some error:\n\nawesome!",
			want: "some error:<br><br>awesome!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFlashErrorString(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}
