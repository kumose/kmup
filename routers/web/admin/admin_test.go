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

package admin

import (
	"net/http"
	"testing"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestShadowPassword(t *testing.T) {
	kases := []struct {
		Provider string
		CfgItem  string
		Result   string
	}{
		{
			Provider: "redis",
			CfgItem:  "network=tcp,addr=:6379,password=kmup,db=0,pool_size=100,idle_timeout=180",
			Result:   "network=tcp,addr=:6379,password=******,db=0,pool_size=100,idle_timeout=180",
		},
		{
			Provider: "mysql",
			CfgItem:  "root:@tcp(localhost:3306)/kmup?charset=utf8",
			Result:   "root:******@tcp(localhost:3306)/kmup?charset=utf8",
		},
		{
			Provider: "mysql",
			CfgItem:  "/kmup?charset=utf8",
			Result:   "/kmup?charset=utf8",
		},
		{
			Provider: "mysql",
			CfgItem:  "user:mypassword@/dbname",
			Result:   "user:******@/dbname",
		},
		{
			Provider: "postgres",
			CfgItem:  "user=pqgotest dbname=pqgotest sslmode=verify-full",
			Result:   "user=pqgotest dbname=pqgotest sslmode=verify-full",
		},
		{
			Provider: "postgres",
			CfgItem:  "user=pqgotest password= dbname=pqgotest sslmode=verify-full",
			Result:   "user=pqgotest password=****** dbname=pqgotest sslmode=verify-full",
		},
		{
			Provider: "postgres",
			CfgItem:  "postgres://user:pass@hostname/dbname",
			Result:   "postgres://user:******@hostname/dbname",
		},
		{
			Provider: "couchbase",
			CfgItem:  "http://dev-couchbase.example.com:8091/",
			Result:   "http://dev-couchbase.example.com:8091/",
		},
		{
			Provider: "couchbase",
			CfgItem:  "http://user:the_password@dev-couchbase.example.com:8091/",
			Result:   "http://user:******@dev-couchbase.example.com:8091/",
		},
	}

	for _, k := range kases {
		assert.Equal(t, k.Result, shadowPassword(k.Provider, k.CfgItem))
	}
}

func TestSelfCheckPost(t *testing.T) {
	defer test.MockVariableValue(&setting.AppURL, "http://config/sub/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()

	ctx, resp := contexttest.MockContext(t, "GET http://host/sub/admin/self_check?location_origin=http://frontend")
	SelfCheckPost(ctx)
	assert.Equal(t, http.StatusOK, resp.Code)

	data := struct {
		Problems []string `json:"problems"`
	}{}
	err := json.Unmarshal(resp.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		ctx.Locale.TrString("admin.self_check.location_origin_mismatch", "http://frontend/sub/", "http://config/sub/"),
	}, data.Problems)
}
