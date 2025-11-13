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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAuthorizationHeader(t *testing.T) {
	type parsed = ParsedAuthorizationHeader
	type basic = BasicAuth
	type bearer = BearerToken
	cases := []struct {
		headerValue string
		expected    parsed
		ok          bool
	}{
		{"", parsed{}, false},
		{"?", parsed{}, false},
		{"foo", parsed{}, false},
		{"any value", parsed{}, false},

		{"Basic ?", parsed{}, false},
		{"Basic " + base64.StdEncoding.EncodeToString([]byte("foo")), parsed{}, false},
		{"Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar")), parsed{BasicAuth: &basic{"foo", "bar"}}, true},
		{"basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar")), parsed{BasicAuth: &basic{"foo", "bar"}}, true},

		{"token value", parsed{BearerToken: &bearer{"value"}}, true},
		{"Token value", parsed{BearerToken: &bearer{"value"}}, true},
		{"bearer value", parsed{BearerToken: &bearer{"value"}}, true},
		{"Bearer value", parsed{BearerToken: &bearer{"value"}}, true},
		{"Bearer wrong value", parsed{}, false},
	}
	for _, c := range cases {
		ret, ok := ParseAuthorizationHeader(c.headerValue)
		assert.Equal(t, c.ok, ok, "header %q", c.headerValue)
		assert.Equal(t, c.expected, ret, "header %q", c.headerValue)
	}
}
