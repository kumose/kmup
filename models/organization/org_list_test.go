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

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/structs"

	"github.com/stretchr/testify/assert"
)

func TestOrgList(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	t.Run("CountOrganizations", testCountOrganizations)
	t.Run("FindOrgs", testFindOrgs)
	t.Run("GetUserOrgsList", testGetUserOrgsList)
	t.Run("LoadOrgListTeams", testLoadOrgListTeams)
	t.Run("DoerViewOtherVisibility", testDoerViewOtherVisibility)
}

func testCountOrganizations(t *testing.T) {
	expected, err := db.GetEngine(t.Context()).Where("type=?", user_model.UserTypeOrganization).Count(&organization.Organization{})
	assert.NoError(t, err)
	cnt, err := db.Count[organization.Organization](t.Context(), organization.FindOrgOptions{IncludeVisibility: structs.VisibleTypePrivate})
	assert.NoError(t, err)
	assert.Equal(t, expected, cnt)
}

func testFindOrgs(t *testing.T) {
	orgs, err := db.Find[organization.Organization](t.Context(), organization.FindOrgOptions{
		UserID:            4,
		IncludeVisibility: structs.VisibleTypePrivate,
	})
	assert.NoError(t, err)
	if assert.Len(t, orgs, 1) {
		assert.EqualValues(t, 3, orgs[0].ID)
	}

	orgs, err = db.Find[organization.Organization](t.Context(), organization.FindOrgOptions{
		UserID: 4,
	})
	assert.NoError(t, err)
	assert.Empty(t, orgs)

	total, err := db.Count[organization.Organization](t.Context(), organization.FindOrgOptions{
		UserID:            4,
		IncludeVisibility: structs.VisibleTypePrivate,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, total)
}

func testGetUserOrgsList(t *testing.T) {
	orgs, err := organization.GetUserOrgsList(t.Context(), &user_model.User{ID: 4})
	assert.NoError(t, err)
	if assert.Len(t, orgs, 1) {
		assert.EqualValues(t, 3, orgs[0].ID)
		// repo_id: 3 is in the team, 32 is public, 5 is private with no team
		assert.Equal(t, 2, orgs[0].NumRepos)
	}
}

func testLoadOrgListTeams(t *testing.T) {
	orgs, err := organization.GetUserOrgsList(t.Context(), &user_model.User{ID: 4})
	assert.NoError(t, err)
	assert.Len(t, orgs, 1)
	teamsMap, err := organization.OrgList(orgs).LoadTeams(t.Context())
	assert.NoError(t, err)
	assert.Len(t, teamsMap, 1)
	assert.Len(t, teamsMap[3], 5)
}

func testDoerViewOtherVisibility(t *testing.T) {
	assert.Equal(t, structs.VisibleTypePublic, organization.DoerViewOtherVisibility(nil, nil))
	assert.Equal(t, structs.VisibleTypeLimited, organization.DoerViewOtherVisibility(&user_model.User{ID: 1}, &user_model.User{ID: 2}))
	assert.Equal(t, structs.VisibleTypePrivate, organization.DoerViewOtherVisibility(&user_model.User{ID: 1}, &user_model.User{ID: 1}))
	assert.Equal(t, structs.VisibleTypePrivate, organization.DoerViewOtherVisibility(&user_model.User{ID: 1, IsAdmin: true}, &user_model.User{ID: 2}))
}
