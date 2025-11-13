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
	"context"

	"github.com/go-git/go-git/v5/plumbing"
	cgobject "github.com/go-git/go-git/v5/plumbing/object/commitgraph"
)

// CacheCommit will cache the commit from the gitRepository
func (c *Commit) CacheCommit(ctx context.Context) error {
	if c.repo.LastCommitCache == nil {
		return nil
	}
	commitNodeIndex, _ := c.repo.CommitNodeIndex()

	index, err := commitNodeIndex.Get(plumbing.Hash(c.ID.RawValue()))
	if err != nil {
		return err
	}

	return c.recursiveCache(ctx, index, &c.Tree, "", 1)
}

func (c *Commit) recursiveCache(ctx context.Context, index cgobject.CommitNode, tree *Tree, treePath string, level int) error {
	if level == 0 {
		return nil
	}

	entries, err := tree.ListEntries()
	if err != nil {
		return err
	}

	entryPaths := make([]string, len(entries))
	entryMap := make(map[string]*TreeEntry)
	for i, entry := range entries {
		entryPaths[i] = entry.Name()
		entryMap[entry.Name()] = entry
	}

	commits, err := GetLastCommitForPaths(ctx, c.repo.LastCommitCache, index, treePath, entryPaths)
	if err != nil {
		return err
	}

	for entry := range commits {
		if entryMap[entry].IsDir() {
			subTree, err := tree.SubTree(entry)
			if err != nil {
				return err
			}
			if err := c.recursiveCache(ctx, index, subTree, entry, level-1); err != nil {
				return err
			}
		}
	}

	return nil
}
