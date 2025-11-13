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

	"github.com/stretchr/testify/assert"
)

func TestCommitSubmoduleLink(t *testing.T) {
	assert.Nil(t, (*CommitSubmoduleFile)(nil).SubmoduleWebLinkTree(t.Context()))
	assert.Nil(t, (*CommitSubmoduleFile)(nil).SubmoduleWebLinkCompare(t.Context(), "", ""))
	assert.Nil(t, (&CommitSubmoduleFile{}).SubmoduleWebLinkTree(t.Context()))
	assert.Nil(t, (&CommitSubmoduleFile{}).SubmoduleWebLinkCompare(t.Context(), "", ""))

	t.Run("GitHubRepo", func(t *testing.T) {
		sf := NewCommitSubmoduleFile("/any/repo-link", "full-path", "git@github.com:user/repo.git", "aaaa")
		wl := sf.SubmoduleWebLinkTree(t.Context())
		assert.Equal(t, "https://github.com/user/repo", wl.RepoWebLink)
		assert.Equal(t, "https://github.com/user/repo/tree/aaaa", wl.CommitWebLink)

		wl = sf.SubmoduleWebLinkCompare(t.Context(), "1111", "2222")
		assert.Equal(t, "https://github.com/user/repo", wl.RepoWebLink)
		assert.Equal(t, "https://github.com/user/repo/compare/1111...2222", wl.CommitWebLink)
	})

	t.Run("RelativePath", func(t *testing.T) {
		sf := NewCommitSubmoduleFile("/subpath/any/repo-home-link", "full-path", "../../user/repo", "aaaa")
		wl := sf.SubmoduleWebLinkTree(t.Context())
		assert.Equal(t, "/subpath/user/repo", wl.RepoWebLink)
		assert.Equal(t, "/subpath/user/repo/tree/aaaa", wl.CommitWebLink)

		sf = NewCommitSubmoduleFile("/subpath/any/repo-home-link", "dir/submodule", "../../user/repo", "aaaa")
		wl = sf.SubmoduleWebLinkCompare(t.Context(), "1111", "2222")
		assert.Equal(t, "/subpath/user/repo", wl.RepoWebLink)
		assert.Equal(t, "/subpath/user/repo/compare/1111...2222", wl.CommitWebLink)
	})
}
