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
	"crypto/subtle"
	"net/http"

	"github.com/kumose/kmup/modules/setting"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics validate auth token and render prometheus metrics
func Metrics(resp http.ResponseWriter, req *http.Request) {
	if setting.Metrics.Token == "" {
		promhttp.Handler().ServeHTTP(resp, req)
		return
	}
	header := req.Header.Get("Authorization")
	if header == "" {
		http.Error(resp, "", http.StatusUnauthorized)
		return
	}
	got := []byte(header)
	want := []byte("Bearer " + setting.Metrics.Token)
	if subtle.ConstantTimeCompare(got, want) != 1 {
		http.Error(resp, "", http.StatusUnauthorized)
		return
	}
	promhttp.Handler().ServeHTTP(resp, req)
}
