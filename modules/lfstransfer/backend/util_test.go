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

package backend

import (
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestToInternalLFSURL(t *testing.T) {
	defer test.MockVariableValue(&setting.LocalURL, "http://localurl/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()
	cases := []struct {
		url      string
		expected string
	}{
		{"http://appurl/any", ""},
		{"http://appurl/sub/any", ""},
		{"http://appurl/sub/owner/repo/any", ""},
		{"http://appurl/sub/owner/repo/info/any", ""},
		{"http://appurl/sub/owner/repo/info/lfs/any", "http://localurl/api/internal/repo/owner/repo/info/lfs/any"},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, toInternalLFSURL(c.url), c.url)
	}
}

func TestIsInternalLFSURL(t *testing.T) {
	defer test.MockVariableValue(&setting.LocalURL, "http://localurl/")()
	defer test.MockVariableValue(&setting.InternalToken, "mock-token")()
	cases := []struct {
		url      string
		expected bool
	}{
		{"", false},
		{"http://otherurl/api/internal/repo/owner/repo/info/lfs/any", false},
		{"http://localurl/api/internal/repo/owner/repo/info/lfs/any", true},
		{"http://localurl/api/internal/repo/owner/repo/info", false},
		{"http://localurl/api/internal/misc/owner/repo/info/lfs/any", false},
		{"http://localurl/api/internal/owner/repo/info/lfs/any", false},
		{"http://localurl/api/internal/foo/bar", false},
	}
	for _, c := range cases {
		req := newInternalRequestLFS(t.Context(), c.url, "GET", nil, nil)
		assert.Equal(t, c.expected, req != nil, c.url)
		assert.Equal(t, c.expected, isInternalLFSURL(c.url), c.url)
	}
}
