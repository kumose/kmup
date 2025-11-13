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

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestGoGet(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	req := NewRequest(t, "GET", "/blah/glah/plah?go-get=1")
	resp := MakeRequest(t, req, http.StatusOK)

	expected := fmt.Sprintf(`<!doctype html>
<html>
	<head>
		<meta name="go-import" content="%[1]s:%[2]s/blah/glah git %[3]sblah/glah.git">
		<meta name="go-source" content="%[1]s:%[2]s/blah/glah _ %[3]sblah/glah/src/branch/master{/dir} %[3]sblah/glah/src/branch/master{/dir}/{file}#L{line}">
	</head>
	<body>
		go get --insecure %[1]s:%[2]s/blah/glah
	</body>
</html>`, setting.Domain, setting.HTTPPort, setting.AppURL)

	assert.Equal(t, expected, resp.Body.String())
}

func TestGoGetForSSH(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	old := setting.Repository.GoGetCloneURLProtocol
	defer func() {
		setting.Repository.GoGetCloneURLProtocol = old
	}()
	setting.Repository.GoGetCloneURLProtocol = "ssh"

	req := NewRequest(t, "GET", "/blah/glah/plah?go-get=1")
	resp := MakeRequest(t, req, http.StatusOK)

	expected := fmt.Sprintf(`<!doctype html>
<html>
	<head>
		<meta name="go-import" content="%[1]s:%[2]s/blah/glah git ssh://git@%[4]s:%[5]d/blah/glah.git">
		<meta name="go-source" content="%[1]s:%[2]s/blah/glah _ %[3]sblah/glah/src/branch/master{/dir} %[3]sblah/glah/src/branch/master{/dir}/{file}#L{line}">
	</head>
	<body>
		go get --insecure %[1]s:%[2]s/blah/glah
	</body>
</html>`, setting.Domain, setting.HTTPPort, setting.AppURL, setting.SSH.Domain, setting.SSH.Port)

	assert.Equal(t, expected, resp.Body.String())
}
