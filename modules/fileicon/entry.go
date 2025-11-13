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

package fileicon

import "github.com/kumose/kmup/modules/git"

type EntryInfo struct {
	BaseName      string
	EntryMode     git.EntryMode
	SymlinkToMode git.EntryMode
	IsOpen        bool
}

func EntryInfoFromGitTreeEntry(commit *git.Commit, fullPath string, gitEntry *git.TreeEntry) *EntryInfo {
	ret := &EntryInfo{BaseName: gitEntry.Name(), EntryMode: gitEntry.Mode()}
	if gitEntry.IsLink() {
		if res, err := git.EntryFollowLink(commit, fullPath, gitEntry); err == nil && res.TargetEntry.IsDir() {
			ret.SymlinkToMode = res.TargetEntry.Mode()
		}
	}
	return ret
}

func EntryInfoFolder() *EntryInfo {
	return &EntryInfo{EntryMode: git.EntryModeTree}
}

func EntryInfoFolderOpen() *EntryInfo {
	return &EntryInfo{EntryMode: git.EntryModeTree, IsOpen: true}
}
