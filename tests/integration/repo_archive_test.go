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
	"io"
	"net/http"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/routers"
	"github.com/kumose/kmup/routers/web"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestRepoDownloadArchive(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.EnableGzip, true)()
	defer test.MockVariableValue(&web.GzipMinSize, 10)()
	defer test.MockVariableValue(&testWebRoutes, routers.NormalRoutes())()

	req := NewRequest(t, "GET", "/user2/repo1/archive/master.zip")
	req.Header.Set("Accept-Encoding", "gzip")
	resp := MakeRequest(t, req, http.StatusOK)
	bs, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Empty(t, resp.Header().Get("Content-Encoding"))
	assert.Len(t, bs, 320)
}
