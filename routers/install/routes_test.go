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

package install

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	r := Routes()
	assert.NotNil(t, r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `class="page-content install"`)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/no-such", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/assets/img/kmup.svg", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
