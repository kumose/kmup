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

package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsParentDirectorySeparator(t *testing.T) {
	tests := []struct {
		v string
		b bool
	}{
		{
			v: `user2/repo1/info/refs`,
			b: false,
		},
		{
			v: `user2/repo1/HEAD`,
			b: false,
		},
		{
			v: `user2/repo1/some.../strange_file...mp3`,
			b: false,
		},
		{
			v: `user2/repo1/../../custom/conf/app.ini`,
			b: true,
		},
		{
			v: `user2/repo1/objects/info/..\..\..\..\custom\conf\app.ini`,
			b: true,
		},
	}

	for i := range tests {
		assert.Equal(t, tests[i].b, containsParentDirectorySeparator(tests[i].v))
	}
}
