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

package common

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/reqctx"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestRenderPanicErrorPage(t *testing.T) {
	w := httptest.NewRecorder()
	req := &http.Request{URL: &url.URL{}}
	req = req.WithContext(reqctx.NewRequestContextForTest(t.Context()))
	RenderPanicErrorPage(w, req, errors.New("fake panic error (for test only)"))
	respContent := w.Body.String()
	assert.Contains(t, respContent, `class="page-content status-page-500"`)
	assert.Contains(t, respContent, `</html>`)
	assert.Contains(t, respContent, `lang="en-US"`) // make sure the locale work

	// the 500 page doesn't have normal pages footer, it makes it easier to distinguish a normal page and a failed page.
	// especially when a sub-template causes page error, the HTTP response code is still 200,
	// the different "footer" is the only way to know whether a page is fully rendered without error.
	assert.False(t, test.IsNormalPageCompleted(respContent))
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
