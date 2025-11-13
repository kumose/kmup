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

	auth_model "github.com/kumose/kmup/models/auth"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIRepoTeams(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// publicOrgRepo = org3/repo21
	publicOrgRepo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 32})
	// user4
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

	// ListTeams
	req := NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/teams", publicOrgRepo.FullName())).
		AddTokenAuth(token)
	res := MakeRequest(t, req, http.StatusOK)
	var teams []*api.Team
	DecodeJSON(t, res, &teams)
	if assert.Len(t, teams, 2) {
		assert.Equal(t, "Owners", teams[0].Name)
		assert.True(t, teams[0].CanCreateOrgRepo)
		assert.True(t, util.SliceSortedEqual(unit.AllUnitKeyNames(), teams[0].Units), "%v == %v", unit.AllUnitKeyNames(), teams[0].Units)
		assert.Equal(t, "owner", teams[0].Permission)

		assert.Equal(t, "test_team", teams[1].Name)
		assert.False(t, teams[1].CanCreateOrgRepo)
		assert.Equal(t, []string{"repo.issues"}, teams[1].Units)
		assert.Equal(t, "write", teams[1].Permission)
	}

	// IsTeam
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/teams/%s", publicOrgRepo.FullName(), "Test_Team")).
		AddTokenAuth(token)
	res = MakeRequest(t, req, http.StatusOK)
	var team *api.Team
	DecodeJSON(t, res, &team)
	assert.Equal(t, teams[1], team)

	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/teams/%s", publicOrgRepo.FullName(), "NonExistingTeam")).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNotFound)

	// AddTeam with user4
	req = NewRequest(t, "PUT", fmt.Sprintf("/api/v1/repos/%s/teams/%s", publicOrgRepo.FullName(), "team1")).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusForbidden)

	// AddTeam with user2
	user = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	session = loginUser(t, user.Name)
	token = getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
	req = NewRequest(t, "PUT", fmt.Sprintf("/api/v1/repos/%s/teams/%s", publicOrgRepo.FullName(), "team1")).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)
	MakeRequest(t, req, http.StatusUnprocessableEntity) // test duplicate request

	// DeleteTeam
	req = NewRequest(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/teams/%s", publicOrgRepo.FullName(), "team1")).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)
	MakeRequest(t, req, http.StatusUnprocessableEntity) // test duplicate request
}
