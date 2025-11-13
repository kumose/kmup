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
	"net/url"
	"strconv"
	"strings"

	"github.com/kumose/kmup/modules/setting"

	"github.com/go-chi/chi/v5"
)

// PathParam returns the param in request path, eg: "/{var}" => "/a%2fb", then `var == "a/b"`
func (b *Base) PathParam(name string) string {
	s, err := url.PathUnescape(b.PathParamRaw(name))
	if err != nil && !setting.IsProd {
		panic("Failed to unescape path param: " + err.Error() + ", there seems to be a double-unescaping bug")
	}
	return s
}

// PathParamRaw returns the raw param in request path, eg: "/{var}" => "/a%2fb", then `var == "a%2fb"`
func (b *Base) PathParamRaw(name string) string {
	if strings.HasPrefix(name, ":") {
		setting.PanicInDevOrTesting("path param should not start with ':'")
		name = name[1:]
	}
	return chi.URLParam(b.Req, name)
}

// PathParamInt64 returns the param in request path as int64
func (b *Base) PathParamInt64(p string) int64 {
	v, _ := strconv.ParseInt(b.PathParam(p), 10, 64)
	return v
}

func (b *Base) PathParamInt(p string) int {
	v, _ := strconv.Atoi(b.PathParam(p))
	return v
}

// SetPathParam set request path params into routes
func (b *Base) SetPathParam(name, value string) {
	if strings.HasPrefix(name, ":") {
		setting.PanicInDevOrTesting("path param should not start with ':'")
		name = name[1:]
	}
	chi.RouteContext(b).URLParams.Add(name, url.PathEscape(value))
}
