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

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestDownloadByID(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")

	// Request raw blob
	req := NewRequest(t, "GET", "/user2/repo1/raw/blob/4b4851ad51df6a7d9f25c979345979eaeb5b349f")
	resp := session.MakeRequest(t, req, http.StatusOK)

	assert.Equal(t, "# repo1\n\nDescription for repo1", resp.Body.String())
}

func TestDownloadByIDForSVGUsesSecureHeaders(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")

	// Request raw blob
	req := NewRequest(t, "GET", "/user2/repo2/raw/blob/6395b68e1feebb1e4c657b4f9f6ba2676a283c0b")
	resp := session.MakeRequest(t, req, http.StatusOK)

	assert.Equal(t, "default-src 'none'; style-src 'unsafe-inline'; sandbox", resp.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "image/svg+xml", resp.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", resp.Header().Get("X-Content-Type-Options"))
}

func TestDownloadByIDMedia(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")

	// Request raw blob
	req := NewRequest(t, "GET", "/user2/repo1/media/blob/4b4851ad51df6a7d9f25c979345979eaeb5b349f")
	resp := session.MakeRequest(t, req, http.StatusOK)

	assert.Equal(t, "# repo1\n\nDescription for repo1", resp.Body.String())
}

func TestDownloadByIDMediaForSVGUsesSecureHeaders(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")

	// Request raw blob
	req := NewRequest(t, "GET", "/user2/repo2/media/blob/6395b68e1feebb1e4c657b4f9f6ba2676a283c0b")
	resp := session.MakeRequest(t, req, http.StatusOK)

	assert.Equal(t, "default-src 'none'; style-src 'unsafe-inline'; sandbox", resp.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "image/svg+xml", resp.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", resp.Header().Get("X-Content-Type-Options"))
}

func TestDownloadRawTextFileWithoutMimeTypeMapping(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")

	req := NewRequest(t, "GET", "/user2/repo2/raw/branch/master/test.xml")
	resp := session.MakeRequest(t, req, http.StatusOK)

	assert.Equal(t, "text/plain; charset=utf-8", resp.Header().Get("Content-Type"))
}

func TestDownloadRawTextFileWithMimeTypeMapping(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	setting.MimeTypeMap.Map[".xml"] = "text/xml"
	setting.MimeTypeMap.Enabled = true

	session := loginUser(t, "user2")

	req := NewRequest(t, "GET", "/user2/repo2/raw/branch/master/test.xml")
	resp := session.MakeRequest(t, req, http.StatusOK)

	assert.Equal(t, "text/xml; charset=utf-8", resp.Header().Get("Content-Type"))

	delete(setting.MimeTypeMap.Map, ".xml")
	setting.MimeTypeMap.Enabled = false
}
