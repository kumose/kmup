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
	"net/url"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestRemoveSessionCookieHeader(t *testing.T) {
	w := httptest.NewRecorder()
	w.Header().Add("Set-Cookie", (&http.Cookie{Name: setting.SessionConfig.CookieName, Value: "foo"}).String())
	w.Header().Add("Set-Cookie", (&http.Cookie{Name: "other", Value: "bar"}).String())
	assert.Len(t, w.Header().Values("Set-Cookie"), 2)
	removeSessionCookieHeader(w)
	assert.Len(t, w.Header().Values("Set-Cookie"), 1)
	assert.Contains(t, "other=bar", w.Header().Get("Set-Cookie"))
}

func TestRedirectToCurrentSite(t *testing.T) {
	setting.IsInTesting = true
	defer test.MockVariableValue(&setting.AppURL, "http://localhost:3000/sub/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()
	cases := []struct {
		location string
		want     string
	}{
		{"/", "/sub/"},
		{"http://localhost:3000/sub?k=v", "http://localhost:3000/sub?k=v"},
		{"http://other", "/sub/"},
	}
	for _, c := range cases {
		t.Run(c.location, func(t *testing.T) {
			req := &http.Request{URL: &url.URL{Path: "/"}}
			resp := httptest.NewRecorder()
			base := NewBaseContextForTest(resp, req)
			ctx := NewWebContext(base, nil, nil)
			ctx.RedirectToCurrentSite(c.location)
			redirect := test.RedirectURL(resp)
			assert.Equal(t, c.want, redirect)
		})
	}
}
