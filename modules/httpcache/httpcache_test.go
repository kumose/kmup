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

package httpcache

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func countFormalHeaders(h http.Header) (c int) {
	for k := range h {
		// ignore our headers for internal usage
		if strings.HasPrefix(k, "X-Kmup-") {
			continue
		}
		c++
	}
	return c
}

func TestHandleGenericETagCache(t *testing.T) {
	etag := `"test"`

	t.Run("No_If-None-Match", func(t *testing.T) {
		req := &http.Request{Header: make(http.Header)}
		w := httptest.NewRecorder()

		handled := HandleGenericETagCache(req, w, etag)

		assert.False(t, handled)
		assert.Equal(t, 2, countFormalHeaders(w.Header()))
		assert.Contains(t, w.Header(), "Cache-Control")
		assert.Contains(t, w.Header(), "Etag")
		assert.Equal(t, etag, w.Header().Get("Etag"))
	})
	t.Run("Wrong_If-None-Match", func(t *testing.T) {
		req := &http.Request{Header: make(http.Header)}
		w := httptest.NewRecorder()

		req.Header.Set("If-None-Match", `"wrong etag"`)

		handled := HandleGenericETagCache(req, w, etag)

		assert.False(t, handled)
		assert.Equal(t, 2, countFormalHeaders(w.Header()))
		assert.Contains(t, w.Header(), "Cache-Control")
		assert.Contains(t, w.Header(), "Etag")
		assert.Equal(t, etag, w.Header().Get("Etag"))
	})
	t.Run("Correct_If-None-Match", func(t *testing.T) {
		req := &http.Request{Header: make(http.Header)}
		w := httptest.NewRecorder()

		req.Header.Set("If-None-Match", etag)

		handled := HandleGenericETagCache(req, w, etag)

		assert.True(t, handled)
		assert.Equal(t, 1, countFormalHeaders(w.Header()))
		assert.Contains(t, w.Header(), "Etag")
		assert.Equal(t, etag, w.Header().Get("Etag"))
		assert.Equal(t, http.StatusNotModified, w.Code)
	})
	t.Run("Multiple_Wrong_If-None-Match", func(t *testing.T) {
		req := &http.Request{Header: make(http.Header)}
		w := httptest.NewRecorder()

		req.Header.Set("If-None-Match", `"wrong etag", "wrong etag "`)

		handled := HandleGenericETagCache(req, w, etag)

		assert.False(t, handled)
		assert.Equal(t, 2, countFormalHeaders(w.Header()))
		assert.Contains(t, w.Header(), "Cache-Control")
		assert.Contains(t, w.Header(), "Etag")
		assert.Equal(t, etag, w.Header().Get("Etag"))
	})
	t.Run("Multiple_Correct_If-None-Match", func(t *testing.T) {
		req := &http.Request{Header: make(http.Header)}
		w := httptest.NewRecorder()

		req.Header.Set("If-None-Match", `"wrong etag", `+etag)

		handled := HandleGenericETagCache(req, w, etag)

		assert.True(t, handled)
		assert.Equal(t, 1, countFormalHeaders(w.Header()))
		assert.Contains(t, w.Header(), "Etag")
		assert.Equal(t, etag, w.Header().Get("Etag"))
		assert.Equal(t, http.StatusNotModified, w.Code)
	})
}
