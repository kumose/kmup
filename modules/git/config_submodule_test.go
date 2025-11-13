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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigSubmodule(t *testing.T) {
	input := `
[core]
path = test

[submodule "submodule1"]
	path = path1
	url =	https://kmup.io/foo/foo
	#branch = b1

[other1]
branch = master

[submodule "submodule2"]
	path = path2
	url =	https://kmup.io/bar/bar
	branch = b2

[other2]
branch = master

[submodule "submodule3"]
	path = path3
	url =	https://kmup.io/xxx/xxx
`

	subModules, err := configParseSubModules(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Len(t, subModules.cache, 3)

	sm1, _ := subModules.Get("path1")
	assert.Equal(t, &SubModule{Path: "path1", URL: "https://kmup.io/foo/foo", Branch: ""}, sm1)
	sm2, _ := subModules.Get("path2")
	assert.Equal(t, &SubModule{Path: "path2", URL: "https://kmup.io/bar/bar", Branch: "b2"}, sm2)
	sm3, _ := subModules.Get("path3")
	assert.Equal(t, &SubModule{Path: "path3", URL: "https://kmup.io/xxx/xxx", Branch: ""}, sm3)
}
