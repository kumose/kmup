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

package routing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_shortenFilename(t *testing.T) {
	tests := []struct {
		filename string
		fallback string
		expected string
	}{
		{
			"code.kmup.io/routers/common/logger_context.go",
			"NO_FALLBACK",
			"common/logger_context.go",
		},
		{
			"common/logger_context.go",
			"NO_FALLBACK",
			"common/logger_context.go",
		},
		{
			"logger_context.go",
			"NO_FALLBACK",
			"logger_context.go",
		},
		{
			"",
			"USE_FALLBACK",
			"USE_FALLBACK",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("shortenFilename('%s')", tt.filename), func(t *testing.T) {
			gotShort := shortenFilename(tt.filename, tt.fallback)
			assert.Equal(t, tt.expected, gotShort)
		})
	}
}

func Test_trimAnonymousFunctionSuffix(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"notAnonymous",
			"notAnonymous",
		},
		{
			"anonymous.func1",
			"anonymous",
		},
		{
			"notAnonymous.funca",
			"notAnonymous.funca",
		},
		{
			"anonymous.func100",
			"anonymous",
		},
		{
			"anonymous.func100.func6",
			"anonymous.func100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimAnonymousFunctionSuffix(tt.name)
			assert.Equal(t, tt.want, got)
		})
	}
}
