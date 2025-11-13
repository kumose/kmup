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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockExpensive(t *testing.T) {
	cases := []struct {
		expensive bool
		routePath string
	}{
		{false, "/user/xxx"},
		{false, "/login/xxx"},
		{true, "/{username}/{reponame}/archive/xxx"},
		{true, "/{username}/{reponame}/graph"},
		{true, "/{username}/{reponame}/src/xxx"},
		{true, "/{username}/{reponame}/wiki/xxx"},
		{true, "/{username}/{reponame}/activity/xxx"},
	}
	for _, c := range cases {
		assert.Equal(t, c.expensive, isRoutePathExpensive(c.routePath), "routePath: %s", c.routePath)
	}

	assert.True(t, isRoutePathForLongPolling("/user/events"))
}
