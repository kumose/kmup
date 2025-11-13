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
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitSmartHTTP(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		testGitSmartHTTP(t, u)
		testRenamedRepoRedirect(t)
	})
}

func testGitSmartHTTP(t *testing.T, u *url.URL) {
	kases := []struct {
		method, path string
		code         int
	}{
		{
			path: "user2/repo1/info/refs",
			code: http.StatusOK,
		},
		{
			method: "HEAD",
			path:   "user2/repo1/info/refs",
			code:   http.StatusOK,
		},
		{
			path: "user2/repo1/HEAD",
			code: http.StatusOK,
		},
		{
			path: "user2/repo1/objects/info/alternates",
			code: http.StatusNotFound,
		},
		{
			path: "user2/repo1/objects/info/http-alternates",
			code: http.StatusNotFound,
		},
		{
			path: "user2/repo1/../../custom/conf/app.ini",
			code: http.StatusNotFound,
		},
		{
			path: "user2/repo1/objects/info/../../../../custom/conf/app.ini",
			code: http.StatusNotFound,
		},
		{
			path: `user2/repo1/objects/info/..\..\..\..\custom\conf\app.ini`,
			code: http.StatusBadRequest,
		},
	}

	for _, kase := range kases {
		t.Run(kase.path, func(t *testing.T) {
			req, err := http.NewRequest(util.IfZero(kase.method, "GET"), u.String()+kase.path, nil)
			require.NoError(t, err)
			req.SetBasicAuth("user2", userPassword)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, kase.code, resp.StatusCode)
			_, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
		})
	}
}

func testRenamedRepoRedirect(t *testing.T) {
	defer test.MockVariableValue(&setting.Service.RequireSignInViewStrict, true)()

	// git client requires to get a 301 redirect response before 401 unauthorized response
	req := NewRequest(t, "GET", "/user2/oldrepo1/info/refs")
	resp := MakeRequest(t, req, http.StatusMovedPermanently)
	redirect := resp.Header().Get("Location")
	assert.Equal(t, "/user2/repo1/info/refs", redirect)

	req = NewRequest(t, "GET", redirect)
	resp = MakeRequest(t, req, http.StatusUnauthorized)
	assert.Equal(t, "Unauthorized\n", resp.Body.String())

	req = NewRequest(t, "GET", redirect).AddBasicAuth("user2")
	resp = MakeRequest(t, req, http.StatusOK)
	assert.Contains(t, resp.Body.String(), "65f1bf27bc3bf70f64657658635e66094edbcb4d\trefs/tags/v1.1")
}
