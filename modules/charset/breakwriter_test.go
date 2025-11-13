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

package charset

import (
	"strings"
	"testing"
)

func TestBreakWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		kase    string
		expect  string
		wantErr bool
	}{
		{
			name:   "noline",
			kase:   "abcdefghijklmnopqrstuvwxyz",
			expect: "abcdefghijklmnopqrstuvwxyz",
		},
		{
			name:   "endline",
			kase:   "abcdefghijklmnopqrstuvwxyz\n",
			expect: "abcdefghijklmnopqrstuvwxyz<br>",
		},
		{
			name:   "startline",
			kase:   "\nabcdefghijklmnopqrstuvwxyz",
			expect: "<br>abcdefghijklmnopqrstuvwxyz",
		},
		{
			name:   "onlyline",
			kase:   "\n\n\n",
			expect: "<br><br><br>",
		},
		{
			name:   "empty",
			kase:   "",
			expect: "",
		},
		{
			name:   "midline",
			kase:   "\nabc\ndefghijkl\nmnopqrstuvwxy\nz",
			expect: "<br>abc<br>defghijkl<br>mnopqrstuvwxy<br>z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &strings.Builder{}
			b := &BreakWriter{
				Writer: buf,
			}
			n, err := b.Write([]byte(tt.kase))
			if (err != nil) != tt.wantErr {
				t.Errorf("BreakWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if n != len(tt.kase) {
				t.Errorf("BreakWriter.Write() = %v, want %v", n, len(tt.kase))
			}
			if buf.String() != tt.expect {
				t.Errorf("BreakWriter.Write() wrote %q, want %v", buf.String(), tt.expect)
			}
		})
	}
}
