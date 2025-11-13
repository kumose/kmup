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

package httplib

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeContentByReader(t *testing.T) {
	data := "0123456789abcdef"

	test := func(t *testing.T, expectedStatusCode int, expectedContent string) {
		_, rangeStr, _ := strings.Cut(t.Name(), "_range_")
		r := &http.Request{Header: http.Header{}, Form: url.Values{}}
		if rangeStr != "" {
			r.Header.Set("Range", "bytes="+rangeStr)
		}
		reader := strings.NewReader(data)
		w := httptest.NewRecorder()
		ServeContentByReader(r, w, int64(len(data)), reader, &ServeHeaderOptions{})
		assert.Equal(t, expectedStatusCode, w.Code)
		if expectedStatusCode == http.StatusPartialContent || expectedStatusCode == http.StatusOK {
			assert.Equal(t, strconv.Itoa(len(expectedContent)), w.Header().Get("Content-Length"))
			assert.Equal(t, expectedContent, w.Body.String())
		}
	}

	t.Run("_range_", func(t *testing.T) {
		test(t, http.StatusOK, data)
	})
	t.Run("_range_0-", func(t *testing.T) {
		test(t, http.StatusPartialContent, data)
	})
	t.Run("_range_0-15", func(t *testing.T) {
		test(t, http.StatusPartialContent, data)
	})
	t.Run("_range_1-", func(t *testing.T) {
		test(t, http.StatusPartialContent, data[1:])
	})
	t.Run("_range_1-3", func(t *testing.T) {
		test(t, http.StatusPartialContent, data[1:3+1])
	})
	t.Run("_range_16-", func(t *testing.T) {
		test(t, http.StatusRequestedRangeNotSatisfiable, "")
	})
	t.Run("_range_1-99999", func(t *testing.T) {
		test(t, http.StatusPartialContent, data[1:])
	})
}

func TestServeContentByReadSeeker(t *testing.T) {
	data := "0123456789abcdef"
	tmpFile := t.TempDir() + "/test"
	err := os.WriteFile(tmpFile, []byte(data), 0o644)
	assert.NoError(t, err)

	test := func(t *testing.T, expectedStatusCode int, expectedContent string) {
		_, rangeStr, _ := strings.Cut(t.Name(), "_range_")
		r := &http.Request{Header: http.Header{}, Form: url.Values{}}
		if rangeStr != "" {
			r.Header.Set("Range", "bytes="+rangeStr)
		}

		seekReader, err := os.OpenFile(tmpFile, os.O_RDONLY, 0o644)
		require.NoError(t, err)
		defer seekReader.Close()

		w := httptest.NewRecorder()
		ServeContentByReadSeeker(r, w, nil, seekReader, &ServeHeaderOptions{})
		assert.Equal(t, expectedStatusCode, w.Code)
		if expectedStatusCode == http.StatusPartialContent || expectedStatusCode == http.StatusOK {
			assert.Equal(t, strconv.Itoa(len(expectedContent)), w.Header().Get("Content-Length"))
			assert.Equal(t, expectedContent, w.Body.String())
		}
	}

	t.Run("_range_", func(t *testing.T) {
		test(t, http.StatusOK, data)
	})
	t.Run("_range_0-", func(t *testing.T) {
		test(t, http.StatusPartialContent, data)
	})
	t.Run("_range_0-15", func(t *testing.T) {
		test(t, http.StatusPartialContent, data)
	})
	t.Run("_range_1-", func(t *testing.T) {
		test(t, http.StatusPartialContent, data[1:])
	})
	t.Run("_range_1-3", func(t *testing.T) {
		test(t, http.StatusPartialContent, data[1:3+1])
	})
	t.Run("_range_16-", func(t *testing.T) {
		test(t, http.StatusRequestedRangeNotSatisfiable, "")
	})
	t.Run("_range_1-99999", func(t *testing.T) {
		test(t, http.StatusPartialContent, data[1:])
	})
}
