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

package fileicon_test

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/fileicon"
	"github.com/kumose/kmup/modules/git"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{FixtureFiles: []string{}})
}

func TestFindIconName(t *testing.T) {
	unittest.PrepareTestEnv(t)
	p := fileicon.DefaultMaterialIconProvider()
	assert.Equal(t, "php", p.FindIconName(&fileicon.EntryInfo{BaseName: "foo.php", EntryMode: git.EntryModeBlob}))
	assert.Equal(t, "php", p.FindIconName(&fileicon.EntryInfo{BaseName: "foo.PHP", EntryMode: git.EntryModeBlob}))
	assert.Equal(t, "javascript", p.FindIconName(&fileicon.EntryInfo{BaseName: "foo.js", EntryMode: git.EntryModeBlob}))
	assert.Equal(t, "visualstudio", p.FindIconName(&fileicon.EntryInfo{BaseName: "foo.vba", EntryMode: git.EntryModeBlob}))
}
