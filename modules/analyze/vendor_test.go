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

package analyze

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsVendor(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"cache/", true},
		{"random/cache/", true},
		{"cache", false},
		{"dependencies/", true},
		{"Dependencies/", true},
		{"dependency/", false},
		{"dist/", true},
		{"dist", false},
		{"random/dist/", true},
		{"random/dist", false},
		{"deps/", true},
		{"configure", true},
		{"a/configure", true},
		{"config.guess", true},
		{"config.guess/", false},
		{".vscode/", true},
		{"doc/_build/", true},
		{"a/docs/_build/", true},
		{"a/dasdocs/_build-vsdoc.js", true},
		{"a/dasdocs/_build-vsdoc.j", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := IsVendor(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}
