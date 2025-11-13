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
	"net/http"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/routers"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestNodeinfo(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.Federation.Enabled, true)()
	defer test.MockVariableValue(&testWebRoutes, routers.NormalRoutes())()

	req := NewRequest(t, "GET", "/api/v1/nodeinfo")
	resp := MakeRequest(t, req, http.StatusOK)
	VerifyJSONSchema(t, resp, "nodeinfo_2.1.json")

	var nodeinfo api.NodeInfo
	DecodeJSON(t, resp, &nodeinfo)
	assert.True(t, nodeinfo.OpenRegistrations)
	assert.Equal(t, "kmup", nodeinfo.Software.Name)
	assert.Equal(t, 29, nodeinfo.Usage.Users.Total)
	assert.Equal(t, 22, nodeinfo.Usage.LocalPosts)
	assert.Equal(t, 3, nodeinfo.Usage.LocalComments)
}
