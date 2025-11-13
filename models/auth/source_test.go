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

package auth_test

import (
	"strings"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm/schemas"
)

type TestSource struct {
	auth_model.ConfigBase

	Provider                      string
	ClientID                      string
	ClientSecret                  string
	OpenIDConnectAutoDiscoveryURL string
	IconURL                       string
}

// FromDB fills up a LDAPConfig from serialized format.
func (source *TestSource) FromDB(bs []byte) error {
	return json.Unmarshal(bs, &source)
}

// ToDB exports a LDAPConfig to a serialized format.
func (source *TestSource) ToDB() ([]byte, error) {
	return json.Marshal(source)
}

func TestDumpAuthSource(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	authSourceSchema, err := unittest.GetXORMEngine().TableInfo(new(auth_model.Source))
	assert.NoError(t, err)

	auth_model.RegisterTypeConfig(auth_model.OAuth2, new(TestSource))

	auth_model.CreateSource(t.Context(), &auth_model.Source{
		Type:     auth_model.OAuth2,
		Name:     "TestSource",
		IsActive: false,
		Cfg: &TestSource{
			Provider: "ConvertibleSourceName",
			ClientID: "42",
		},
	})

	sb := new(strings.Builder)

	// TODO: this test is quite hacky, it should use a low-level "select" (without model processors) but not a database dump
	engine := unittest.GetXORMEngine()
	require.NoError(t, engine.DumpTables([]*schemas.Table{authSourceSchema}, sb))
	assert.Contains(t, sb.String(), `"Provider":"ConvertibleSourceName"`)
}
