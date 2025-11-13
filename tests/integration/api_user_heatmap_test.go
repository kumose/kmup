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
	"time"

	activities_model "github.com/kumose/kmup/models/activities"
	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestUserHeatmap(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	adminUsername := "user1"
	normalUsername := "user2"
	token := getUserToken(t, adminUsername, auth_model.AccessTokenScopeReadUser)

	fakeNow := time.Date(2011, 10, 20, 0, 0, 0, 0, time.Local)
	timeutil.MockSet(fakeNow)
	defer timeutil.MockUnset()

	req := NewRequest(t, "GET", fmt.Sprintf("/api/v1/users/%s/heatmap", normalUsername)).
		AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusOK)
	var heatmap []*activities_model.UserHeatmapData
	DecodeJSON(t, resp, &heatmap)
	var dummyheatmap []*activities_model.UserHeatmapData
	dummyheatmap = append(dummyheatmap, &activities_model.UserHeatmapData{Timestamp: 1603227600, Contributions: 1})

	assert.Equal(t, dummyheatmap, heatmap)
}
