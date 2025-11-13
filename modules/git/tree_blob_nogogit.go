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

//go:build !gogit

package git

import (
	"path"
	"strings"
)

// GetTreeEntryByPath get the tree entries according the sub dir
func (t *Tree) GetTreeEntryByPath(relpath string) (_ *TreeEntry, err error) {
	if len(relpath) == 0 {
		return &TreeEntry{
			ptree:     t,
			ID:        t.ID,
			name:      "",
			entryMode: EntryModeTree,
		}, nil
	}

	relpath = path.Clean(relpath)
	parts := strings.Split(relpath, "/")

	tree := t
	for _, name := range parts[:len(parts)-1] {
		tree, err = tree.SubTree(name)
		if err != nil {
			return nil, err
		}
	}

	name := parts[len(parts)-1]
	entries, err := tree.ListEntries()
	if err != nil {
		return nil, err
	}
	for _, v := range entries {
		if v.Name() == name {
			return v, nil
		}
	}
	return nil, ErrNotExist{"", relpath}
}
