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

package organization_test

import (
	"testing"

	org_model "github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func Test_GetTeamsByIDs(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// 1 owner team, 2 normal team
	teams, err := org_model.GetTeamsByIDs(t.Context(), []int64{1, 2})
	assert.NoError(t, err)
	assert.Len(t, teams, 2)
	assert.Equal(t, "Owners", teams[1].Name)
	assert.Equal(t, "team1", teams[2].Name)
}
