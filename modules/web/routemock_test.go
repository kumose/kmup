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

package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestRouteMock(t *testing.T) {
	setting.IsInTesting = true

	r := NewRouter()
	middleware1 := func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("X-Test-Middleware1", "m1")
	}
	middleware2 := func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("X-Test-Middleware2", "m2")
	}
	handler := func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("X-Test-Handler", "h")
	}
	r.Get("/foo", middleware1, RouterMockPoint("mock-point"), middleware2, handler)

	// normal request
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/foo", nil)
	assert.NoError(t, err)
	r.ServeHTTP(recorder, req)
	assert.Len(t, recorder.Header(), 3)
	assert.Equal(t, "m1", recorder.Header().Get("X-Test-Middleware1"))
	assert.Equal(t, "m2", recorder.Header().Get("X-Test-Middleware2"))
	assert.Equal(t, "h", recorder.Header().Get("X-Test-Handler"))
	RouteMockReset()

	// mock at "mock-point"
	RouteMock("mock-point", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("X-Test-MockPoint", "a")
		resp.WriteHeader(http.StatusOK)
	})
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "http://localhost:8000/foo", nil)
	assert.NoError(t, err)
	r.ServeHTTP(recorder, req)
	assert.Len(t, recorder.Header(), 2)
	assert.Equal(t, "m1", recorder.Header().Get("X-Test-Middleware1"))
	assert.Equal(t, "a", recorder.Header().Get("X-Test-MockPoint"))
	RouteMockReset()

	// mock at MockAfterMiddlewares
	RouteMock(MockAfterMiddlewares, func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("X-Test-MockPoint", "b")
		resp.WriteHeader(http.StatusOK)
	})
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "http://localhost:8000/foo", nil)
	assert.NoError(t, err)
	r.ServeHTTP(recorder, req)
	assert.Len(t, recorder.Header(), 3)
	assert.Equal(t, "m1", recorder.Header().Get("X-Test-Middleware1"))
	assert.Equal(t, "m2", recorder.Header().Get("X-Test-Middleware2"))
	assert.Equal(t, "b", recorder.Header().Get("X-Test-MockPoint"))
	RouteMockReset()
}
