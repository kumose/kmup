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
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestGenAPILinks(t *testing.T) {
	setting.AppURL = "http://localhost:3326/"
	kases := map[string][]string{
		"api/v1/repos/jerrykan/example-repo/issues?state=all": {
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=2&state=all>; rel="next"`,
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=5&state=all>; rel="last"`,
		},
		"api/v1/repos/jerrykan/example-repo/issues?state=all&page=1": {
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=2&state=all>; rel="next"`,
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=5&state=all>; rel="last"`,
		},
		"api/v1/repos/jerrykan/example-repo/issues?state=all&page=2": {
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=3&state=all>; rel="next"`,
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=5&state=all>; rel="last"`,
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=1&state=all>; rel="first"`,
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=1&state=all>; rel="prev"`,
		},
		"api/v1/repos/jerrykan/example-repo/issues?state=all&page=5": {
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=1&state=all>; rel="first"`,
			`<http://localhost:3326/api/v1/repos/jerrykan/example-repo/issues?page=4&state=all>; rel="prev"`,
		},
	}

	for req, response := range kases {
		u, err := url.Parse(setting.AppURL + req)
		assert.NoError(t, err)

		p := u.Query().Get("page")
		curPage, _ := strconv.Atoi(p)

		links := genAPILinks(u, 100, 20, curPage)

		assert.Equal(t, links, response)
	}
}
