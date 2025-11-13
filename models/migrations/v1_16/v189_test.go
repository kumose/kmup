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

package v1_16

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/json"

	"github.com/stretchr/testify/assert"
)

// LoginSource represents an external way for authorizing users.
type LoginSourceOriginalV189 struct {
	ID        int64 `xorm:"pk autoincr"`
	Type      int
	IsActived bool   `xorm:"INDEX NOT NULL DEFAULT false"`
	Cfg       string `xorm:"TEXT"`
	Expected  string `xorm:"TEXT"`
}

func (ls *LoginSourceOriginalV189) TableName() string {
	return "login_source"
}

func Test_UnwrapLDAPSourceCfg(t *testing.T) {
	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(LoginSourceOriginalV189))
	if x == nil || t.Failed() {
		defer deferable()
		return
	}
	defer deferable()

	// LoginSource represents an external way for authorizing users.
	type LoginSource struct {
		ID       int64 `xorm:"pk autoincr"`
		Type     int
		IsActive bool   `xorm:"INDEX NOT NULL DEFAULT false"`
		Cfg      string `xorm:"TEXT"`
		Expected string `xorm:"TEXT"`
	}

	// Run the migration
	if err := UnwrapLDAPSourceCfg(x); err != nil {
		assert.NoError(t, err)
		return
	}

	const batchSize = 100
	for start := 0; ; start += batchSize {
		sources := make([]*LoginSource, 0, batchSize)
		if err := x.Table("login_source").Limit(batchSize, start).Find(&sources); err != nil {
			assert.NoError(t, err)
			return
		}

		if len(sources) == 0 {
			break
		}

		for _, source := range sources {
			converted := map[string]any{}
			expected := map[string]any{}

			if err := json.Unmarshal([]byte(source.Cfg), &converted); err != nil {
				assert.NoError(t, err)
				return
			}

			if err := json.Unmarshal([]byte(source.Expected), &expected); err != nil {
				assert.NoError(t, err)
				return
			}

			assert.Equal(t, expected, converted, "UnwrapLDAPSourceCfg failed for %d", source.ID)
			assert.Equal(t, source.ID%2 == 0, source.IsActive, "UnwrapLDAPSourceCfg failed for %d", source.ID)
		}
	}
}
