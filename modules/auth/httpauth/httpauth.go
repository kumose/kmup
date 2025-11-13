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

package httpauth

import (
	"encoding/base64"
	"strings"

	"github.com/kumose/kmup/modules/util"
)

type BasicAuth struct {
	Username, Password string
}

type BearerToken struct {
	Token string
}

type ParsedAuthorizationHeader struct {
	BasicAuth   *BasicAuth
	BearerToken *BearerToken
}

func ParseAuthorizationHeader(header string) (ret ParsedAuthorizationHeader, _ bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 {
		return ret, false
	}
	if util.AsciiEqualFold(parts[0], "basic") {
		s, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return ret, false
		}
		u, p, ok := strings.Cut(string(s), ":")
		if !ok {
			return ret, false
		}
		ret.BasicAuth = &BasicAuth{Username: u, Password: p}
		return ret, true
	} else if util.AsciiEqualFold(parts[0], "token") || util.AsciiEqualFold(parts[0], "bearer") {
		ret.BearerToken = &BearerToken{Token: parts[1]}
		return ret, true
	}
	return ret, false
}
