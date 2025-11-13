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
	"testing"

	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFollowLink(t *testing.T) {
	r, err := OpenRepository(t.Context(), "tests/repos/repo1_bare")
	require.NoError(t, err)
	defer r.Close()

	commit, err := r.GetCommit("37991dec2c8e592043f47155ce4808d4580f9123")
	require.NoError(t, err)

	// get the symlink
	{
		lnkFullPath := "foo/bar/link_to_hello"
		lnk, err := commit.Tree.GetTreeEntryByPath("foo/bar/link_to_hello")
		require.NoError(t, err)
		assert.True(t, lnk.IsLink())

		// should be able to dereference to target
		res, err := EntryFollowLink(commit, lnkFullPath, lnk)
		require.NoError(t, err)
		assert.Equal(t, "hello", res.TargetEntry.Name())
		assert.Equal(t, "foo/nar/hello", res.TargetFullPath)
		assert.False(t, res.TargetEntry.IsLink())
		assert.Equal(t, "b14df6442ea5a1b382985a6549b85d435376c351", res.TargetEntry.ID.String())
	}

	{
		// should error when called on a normal file
		entry, err := commit.Tree.GetTreeEntryByPath("file1.txt")
		require.NoError(t, err)
		res, err := EntryFollowLink(commit, "file1.txt", entry)
		assert.ErrorIs(t, err, util.ErrUnprocessableContent)
		assert.Nil(t, res)
	}

	{
		// should error for broken links
		entry, err := commit.Tree.GetTreeEntryByPath("foo/broken_link")
		require.NoError(t, err)
		assert.True(t, entry.IsLink())
		res, err := EntryFollowLink(commit, "foo/broken_link", entry)
		assert.ErrorIs(t, err, util.ErrNotExist)
		assert.Equal(t, "nar/broken_link", res.SymlinkContent)
	}

	{
		// should error for external links
		entry, err := commit.Tree.GetTreeEntryByPath("foo/outside_repo")
		require.NoError(t, err)
		assert.True(t, entry.IsLink())
		res, err := EntryFollowLink(commit, "foo/outside_repo", entry)
		assert.ErrorIs(t, err, util.ErrNotExist)
		assert.Equal(t, "../../outside_repo", res.SymlinkContent)
	}

	{
		// testing fix for short link bug
		entry, err := commit.Tree.GetTreeEntryByPath("foo/link_short")
		require.NoError(t, err)
		res, err := EntryFollowLink(commit, "foo/link_short", entry)
		assert.ErrorIs(t, err, util.ErrNotExist)
		assert.Equal(t, "a", res.SymlinkContent)
	}
}
