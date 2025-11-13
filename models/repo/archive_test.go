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

func TestArchiveType(t *testing.T) {
	name, archiveType := SplitArchiveNameType("test.tar.gz")
	assert.Equal(t, "test", name)
	assert.Equal(t, "tar.gz", archiveType.String())

	name, archiveType = SplitArchiveNameType("a/b/test.zip")
	assert.Equal(t, "a/b/test", name)
	assert.Equal(t, "zip", archiveType.String())

	name, archiveType = SplitArchiveNameType("1234.bundle")
	assert.Equal(t, "1234", name)
	assert.Equal(t, "bundle", archiveType.String())

	name, archiveType = SplitArchiveNameType("test")
	assert.Equal(t, "test", name)
	assert.Equal(t, "unknown", archiveType.String())

	name, archiveType = SplitArchiveNameType("test.xz")
	assert.Equal(t, "test.xz", name)
	assert.Equal(t, "unknown", archiveType.String())
}
