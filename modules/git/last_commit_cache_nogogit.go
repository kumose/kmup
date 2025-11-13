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
	"context"
)

// CacheCommit will cache the commit from the gitRepository
func (c *Commit) CacheCommit(ctx context.Context) error {
	if c.repo.LastCommitCache == nil {
		return nil
	}
	return c.recursiveCache(ctx, &c.Tree, "", 1)
}

func (c *Commit) recursiveCache(ctx context.Context, tree *Tree, treePath string, level int) error {
	if level == 0 {
		return nil
	}

	entries, err := tree.ListEntries()
	if err != nil {
		return err
	}

	entryPaths := make([]string, len(entries))
	for i, entry := range entries {
		entryPaths[i] = entry.Name()
	}

	_, err = WalkGitLog(ctx, c.repo, c, treePath, entryPaths...)
	if err != nil {
		return err
	}

	for _, treeEntry := range entries {
		// entryMap won't contain "" therefore skip this.
		if treeEntry.IsDir() {
			subTree, err := tree.SubTree(treeEntry.Name())
			if err != nil {
				return err
			}
			if err := c.recursiveCache(ctx, subTree, treeEntry.Name(), level-1); err != nil {
				return err
			}
		}
	}

	return nil
}
