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

package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSignatureFromCommitLine(t *testing.T) {
	tests := []struct {
		line string
		want *Signature
	}{
		{
			line: "a b <c@d.com> 12345 +0100",
			want: &Signature{
				Name:  "a b",
				Email: "c@d.com",
				When:  time.Unix(12345, 0).In(time.FixedZone("", 3600)),
			},
		},
		{
			line: "bad line",
			want: &Signature{Name: "bad line"},
		},
		{
			line: "bad < line",
			want: &Signature{Name: "bad < line"},
		},
		{
			line: "bad > line",
			want: &Signature{Name: "bad > line"},
		},
		{
			line: "bad-line <name@example.com>",
			want: &Signature{Name: "bad-line <name@example.com>"},
		},
	}
	for _, test := range tests {
		got := parseSignatureFromCommitLine(test.line)
		assert.Equal(t, test.want, got)
	}
}
