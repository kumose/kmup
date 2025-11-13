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

package org

import (
	"testing"

	"github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestUser_RemoveMember(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	org := unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: 3})
	user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
	user5 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 5})

	// remove a user that is a member
	unittest.AssertExistsAndLoadBean(t, &organization.OrgUser{UID: user4.ID, OrgID: org.ID})
	prevNumMembers := org.NumMembers
	assert.NoError(t, RemoveOrgUser(t.Context(), org, user4))
	unittest.AssertNotExistsBean(t, &organization.OrgUser{UID: user4.ID, OrgID: org.ID})

	org = unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: org.ID})
	assert.Equal(t, prevNumMembers-1, org.NumMembers)

	// remove a user that is not a member
	unittest.AssertNotExistsBean(t, &organization.OrgUser{UID: user5.ID, OrgID: org.ID})
	prevNumMembers = org.NumMembers
	assert.NoError(t, RemoveOrgUser(t.Context(), org, user5))
	unittest.AssertNotExistsBean(t, &organization.OrgUser{UID: user5.ID, OrgID: org.ID})

	org = unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: org.ID})
	assert.Equal(t, prevNumMembers, org.NumMembers)

	unittest.CheckConsistencyFor(t, &user_model.User{}, &organization.Team{})
}

func TestRemoveOrgUser(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	testSuccess := func(org *organization.Organization, user *user_model.User) {
		expectedNumMembers := org.NumMembers
		if unittest.GetBean(t, &organization.OrgUser{OrgID: org.ID, UID: user.ID}) != nil {
			expectedNumMembers--
		}
		assert.NoError(t, RemoveOrgUser(t.Context(), org, user))
		unittest.AssertNotExistsBean(t, &organization.OrgUser{OrgID: org.ID, UID: user.ID})
		org = unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: org.ID})
		assert.Equal(t, expectedNumMembers, org.NumMembers)
	}

	org3 := unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: 3})
	org7 := unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: 7})
	user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
	user5 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 5})

	testSuccess(org3, user4)

	org3 = unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: 3})
	testSuccess(org3, user4)

	err := RemoveOrgUser(t.Context(), org7, user5)
	assert.Error(t, err)
	assert.True(t, organization.IsErrLastOrgOwner(err))
	unittest.AssertExistsAndLoadBean(t, &organization.OrgUser{OrgID: org7.ID, UID: user5.ID})
	unittest.CheckConsistencyFor(t, &user_model.User{}, &organization.Team{})
}
