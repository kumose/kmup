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

// CommitInfo describes the first commit with the provided entry
type CommitInfo struct {
	Entry         *TreeEntry
	Commit        *Commit
	SubmoduleFile *CommitSubmoduleFile
}

func GetCommitInfoSubmoduleFile(repoLink, fullPath string, commit *Commit, refCommitID ObjectID) (*CommitSubmoduleFile, error) {
	submodule, err := commit.GetSubModule(fullPath)
	if err != nil {
		return nil, err
	}
	if submodule == nil {
		// unable to find submodule from ".gitmodules" file
		return NewCommitSubmoduleFile(repoLink, fullPath, "", refCommitID.String()), nil
	}
	return NewCommitSubmoduleFile(repoLink, fullPath, submodule.URL, refCommitID.String()), nil
}
