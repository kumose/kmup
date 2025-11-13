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

//go:build gogit

package git

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/stretchr/testify/assert"
)

func TestEntryGogit(t *testing.T) {
	cases := map[EntryMode]filemode.FileMode{
		EntryModeBlob:    filemode.Regular,
		EntryModeCommit:  filemode.Submodule,
		EntryModeExec:    filemode.Executable,
		EntryModeSymlink: filemode.Symlink,
		EntryModeTree:    filemode.Dir,
	}
	for emode, fmode := range cases {
		assert.EqualValues(t, fmode, entryModeToGogitFileMode(emode))
		assert.EqualValues(t, emode, gogitFileModeToEntryMode(fmode))
	}
}
