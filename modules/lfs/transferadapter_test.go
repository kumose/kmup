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

package lfs

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/json"

	"github.com/stretchr/testify/assert"
)

func TestBasicTransferAdapterName(t *testing.T) {
	a := &BasicTransferAdapter{}

	assert.Equal(t, "basic", a.Name())
}

func TestBasicTransferAdapter(t *testing.T) {
	p := Pointer{Oid: "b5a2c96250612366ea272ffac6d9744aaf4b45aacd96aa7cfcb931ee3b558259", Size: 5}

	roundTripHandler := func(req *http.Request) *http.Response {
		assert.Equal(t, AcceptHeader, req.Header.Get("Accept"))
		assert.Equal(t, "test-value", req.Header.Get("test-header"))

		url := req.URL.String()
		if strings.Contains(url, "download-request") {
			assert.Equal(t, "GET", req.Method)

			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("dummy"))}
		} else if strings.Contains(url, "upload-request") {
			assert.Equal(t, "PUT", req.Method)
			assert.Equal(t, "application/octet-stream", req.Header.Get("Content-Type"))

			b, err := io.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, "dummy", string(b))

			return &http.Response{StatusCode: http.StatusOK}
		} else if strings.Contains(url, "verify-request") {
			assert.Equal(t, "POST", req.Method)
			assert.Equal(t, MediaType, req.Header.Get("Content-Type"))

			var vp Pointer
			err := json.NewDecoder(req.Body).Decode(&vp)
			assert.NoError(t, err)
			assert.Equal(t, p.Oid, vp.Oid)
			assert.Equal(t, p.Size, vp.Size)

			return &http.Response{StatusCode: http.StatusOK}
		} else if strings.Contains(url, "error-response") {
			er := &ErrorResponse{
				Message: "Object not found",
			}
			payload := new(bytes.Buffer)
			json.NewEncoder(payload).Encode(er)

			return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(payload)}
		}
		t.Errorf("Unknown test case: %s", url)
		return nil
	}

	hc := &http.Client{Transport: RoundTripFunc(roundTripHandler)}
	a := &BasicTransferAdapter{hc}

	t.Run("Download", func(t *testing.T) {
		cases := []struct {
			link          *Link
			expectederror string
		}{
			// case 0
			{
				link: &Link{
					Href:   "https://download-request.io",
					Header: map[string]string{"test-header": "test-value"},
				},
				expectederror: "",
			},
			// case 1
			{
				link: &Link{
					Href:   "https://error-response.io",
					Header: map[string]string{"test-header": "test-value"},
				},
				expectederror: "Object not found",
			},
		}

		for n, c := range cases {
			_, err := a.Download(t.Context(), c.link)
			if len(c.expectederror) > 0 {
				assert.Contains(t, err.Error(), c.expectederror, "case %d: '%s' should contain '%s'", n, err.Error(), c.expectederror)
			} else {
				assert.NoError(t, err, "case %d", n)
			}
		}
	})

	t.Run("Upload", func(t *testing.T) {
		cases := []struct {
			link          *Link
			expectederror string
		}{
			// case 0
			{
				link: &Link{
					Href:   "https://upload-request.io",
					Header: map[string]string{"test-header": "test-value"},
				},
				expectederror: "",
			},
			// case 1
			{
				link: &Link{
					Href:   "https://error-response.io",
					Header: map[string]string{"test-header": "test-value"},
				},
				expectederror: "Object not found",
			},
		}

		for n, c := range cases {
			err := a.Upload(t.Context(), c.link, p, strings.NewReader("dummy"))
			if len(c.expectederror) > 0 {
				assert.Contains(t, err.Error(), c.expectederror, "case %d: '%s' should contain '%s'", n, err.Error(), c.expectederror)
			} else {
				assert.NoError(t, err, "case %d", n)
			}
		}
	})

	t.Run("Verify", func(t *testing.T) {
		cases := []struct {
			link          *Link
			expectederror string
		}{
			// case 0
			{
				link: &Link{
					Href:   "https://verify-request.io",
					Header: map[string]string{"test-header": "test-value"},
				},
				expectederror: "",
			},
			// case 1
			{
				link: &Link{
					Href:   "https://error-response.io",
					Header: map[string]string{"test-header": "test-value"},
				},
				expectederror: "Object not found",
			},
		}

		for n, c := range cases {
			err := a.Verify(t.Context(), c.link, p)
			if len(c.expectederror) > 0 {
				assert.Contains(t, err.Error(), c.expectederror, "case %d: '%s' should contain '%s'", n, err.Error(), c.expectederror)
			} else {
				assert.NoError(t, err, "case %d", n)
			}
		}
	})
}
