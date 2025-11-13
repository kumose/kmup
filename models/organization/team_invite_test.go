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

	"github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestTeamInvite(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	team := unittest.AssertExistsAndLoadBean(t, &organization.Team{ID: 2})

	t.Run("MailExistsInTeam", func(t *testing.T) {
		user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

		// user 2 already added to team 2, should result in error
		_, err := organization.CreateTeamInvite(t.Context(), user2, team, user2.Email)
		assert.Error(t, err)
	})

	t.Run("CreateAndRemove", func(t *testing.T) {
		user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})

		invite, err := organization.CreateTeamInvite(t.Context(), user1, team, "org3@example.com")
		assert.NotNil(t, invite)
		assert.NoError(t, err)

		// Shouldn't allow duplicate invite
		_, err = organization.CreateTeamInvite(t.Context(), user1, team, "org3@example.com")
		assert.Error(t, err)

		// should remove invite
		assert.NoError(t, organization.RemoveInviteByID(t.Context(), invite.ID, invite.TeamID))

		// invite should not exist
		_, err = organization.GetInviteByToken(t.Context(), invite.Token)
		assert.Error(t, err)
	})
}
