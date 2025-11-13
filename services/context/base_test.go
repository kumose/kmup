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

package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	setting.IsInTesting = true
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	cases := []struct {
		url  string
		keep bool
	}{
		{"http://test", false},
		{"https://test", false},
		{"//test", false},
		{"/://test", true},
		{"/test", true},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		b := NewBaseContextForTest(resp, req)
		resp.Header().Add("Set-Cookie", (&http.Cookie{Name: setting.SessionConfig.CookieName, Value: "dummy"}).String())
		b.Redirect(c.url)
		has := resp.Header().Get("Set-Cookie") == "i_like_kmup=dummy"
		assert.Equal(t, c.keep, has, "url = %q", c.url)
	}

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	req.Header.Add("HX-Request", "true")
	b := NewBaseContextForTest(resp, req)
	b.Redirect("/other")
	assert.Equal(t, "/other", resp.Header().Get("HX-Redirect"))
	assert.Equal(t, http.StatusNoContent, resp.Code)
}
