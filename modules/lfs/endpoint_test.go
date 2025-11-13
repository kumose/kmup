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

package lfs

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func str2url(raw string) *url.URL {
	u, _ := url.Parse(raw)
	return u
}

func TestDetermineEndpoint(t *testing.T) {
	// Test cases
	cases := []struct {
		cloneurl string
		lfsurl   string
		expected *url.URL
	}{
		// case 0
		{
			cloneurl: "",
			lfsurl:   "",
			expected: nil,
		},
		// case 1
		{
			cloneurl: "https://git.com/repo",
			lfsurl:   "",
			expected: str2url("https://git.com/repo.git/info/lfs"),
		},
		// case 2
		{
			cloneurl: "https://git.com/repo.git",
			lfsurl:   "",
			expected: str2url("https://git.com/repo.git/info/lfs"),
		},
		// case 3
		{
			cloneurl: "",
			lfsurl:   "https://gitlfs.com/repo",
			expected: str2url("https://gitlfs.com/repo"),
		},
		// case 4
		{
			cloneurl: "https://git.com/repo.git",
			lfsurl:   "https://gitlfs.com/repo",
			expected: str2url("https://gitlfs.com/repo"),
		},
		// case 5
		{
			cloneurl: "git://git.com/repo.git",
			lfsurl:   "",
			expected: str2url("https://git.com/repo.git/info/lfs"),
		},
		// case 6
		{
			cloneurl: "",
			lfsurl:   "git://gitlfs.com/repo",
			expected: str2url("https://gitlfs.com/repo"),
		},
	}

	for n, c := range cases {
		ep := DetermineEndpoint(c.cloneurl, c.lfsurl)

		assert.Equal(t, c.expected, ep, "case %d: error should match", n)
	}
}
