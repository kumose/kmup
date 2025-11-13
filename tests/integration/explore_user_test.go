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
	"net/http"
	"testing"

	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestExploreUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	cases := []struct{ sortOrder, expected string }{
		{"", "?sort=newest&q="},
		{"newest", "?sort=newest&q="},
		{"oldest", "?sort=oldest&q="},
		{"alphabetically", "?sort=alphabetically&q="},
		{"reversealphabetically", "?sort=reversealphabetically&q="},
	}
	for _, c := range cases {
		req := NewRequest(t, "GET", "/explore/users?sort="+c.sortOrder)
		resp := MakeRequest(t, req, http.StatusOK)
		h := NewHTMLParser(t, resp.Body)
		href, _ := h.Find(`.ui.dropdown .menu a.active.item[href^="?sort="]`).Attr("href")
		assert.Equal(t, c.expected, href)
	}

	// these sort orders shouldn't be supported, to avoid leaking user activity
	cases404 := []string{
		"/explore/users?sort=lastlogin",
		"/explore/users?sort=reverselastlogin",
		"/explore/users?sort=leastupdate",
		"/explore/users?sort=reverseleastupdate",
	}
	for _, c := range cases404 {
		req := NewRequest(t, "GET", c).SetHeader("Accept", "text/html")
		MakeRequest(t, req, http.StatusNotFound)
	}
}
