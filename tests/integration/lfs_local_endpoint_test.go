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

package integration

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/kumose/kmup/modules/lfs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func str2url(raw string) *url.URL {
	u, _ := url.Parse(raw)
	return u
}

func TestDetermineLocalEndpoint(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	root := t.TempDir()

	rootdotgit := t.TempDir()
	os.Mkdir(filepath.Join(rootdotgit, ".git"), 0o700)

	lfsroot := t.TempDir()

	// Test cases
	cases := []struct {
		cloneurl string
		lfsurl   string
		expected *url.URL
	}{
		// case 0
		{
			cloneurl: root,
			lfsurl:   "",
			expected: str2url("file://" + root),
		},
		// case 1
		{
			cloneurl: root,
			lfsurl:   lfsroot,
			expected: str2url("file://" + lfsroot),
		},
		// case 2
		{
			cloneurl: "https://git.com/repo.git",
			lfsurl:   lfsroot,
			expected: str2url("file://" + lfsroot),
		},
		// case 3
		{
			cloneurl: rootdotgit,
			lfsurl:   "",
			expected: str2url("file://" + filepath.Join(rootdotgit, ".git")),
		},
		// case 4
		{
			cloneurl: "",
			lfsurl:   rootdotgit,
			expected: str2url("file://" + filepath.Join(rootdotgit, ".git")),
		},
		// case 5
		{
			cloneurl: rootdotgit,
			lfsurl:   rootdotgit,
			expected: str2url("file://" + filepath.Join(rootdotgit, ".git")),
		},
		// case 6
		{
			cloneurl: "file://" + root,
			lfsurl:   "",
			expected: str2url("file://" + root),
		},
		// case 7
		{
			cloneurl: "file://" + root,
			lfsurl:   "file://" + lfsroot,
			expected: str2url("file://" + lfsroot),
		},
		// case 8
		{
			cloneurl: root,
			lfsurl:   "file://" + lfsroot,
			expected: str2url("file://" + lfsroot),
		},
		// case 9
		{
			cloneurl: "",
			lfsurl:   "/does/not/exist",
			expected: nil,
		},
		// case 10
		{
			cloneurl: "",
			lfsurl:   "file:///does/not/exist",
			expected: str2url("file:///does/not/exist"),
		},
	}

	for n, c := range cases {
		ep := lfs.DetermineEndpoint(c.cloneurl, c.lfsurl)

		assert.Equal(t, c.expected, ep, "case %d: error should match", n)
	}
}
