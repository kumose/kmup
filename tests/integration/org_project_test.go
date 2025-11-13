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
	"slices"
	"testing"

	unit_model "github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/tests"
)

func TestOrgProjectAccess(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	disabledRepoUnits := unit_model.DisabledRepoUnitsGet()
	unit_model.DisabledRepoUnitsSet(append(slices.Clone(disabledRepoUnits), unit_model.TypeProjects))
	defer unit_model.DisabledRepoUnitsSet(disabledRepoUnits)

	// repo project, 404
	req := NewRequest(t, "GET", "/user2/repo1/projects")
	MakeRequest(t, req, http.StatusNotFound)

	// user project, 200
	req = NewRequest(t, "GET", "/user2/-/projects")
	MakeRequest(t, req, http.StatusOK)

	// org project, 200
	req = NewRequest(t, "GET", "/org3/-/projects")
	MakeRequest(t, req, http.StatusOK)

	// change the org's visibility to private
	session := loginUser(t, "user2")
	req = NewRequestWithValues(t, "POST", "/org/org3/settings", map[string]string{
		"_csrf":      GetUserCSRFToken(t, session),
		"name":       "org3",
		"visibility": "2",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	// user4 can still access the org's project because its team(team1) has the permission
	session = loginUser(t, "user4")
	req = NewRequest(t, "GET", "/org3/-/projects")
	session.MakeRequest(t, req, http.StatusOK)

	// disable team1's project unit
	session = loginUser(t, "user2")
	req = NewRequestWithValues(t, "POST", "/org/org3/teams/team1/edit", map[string]string{
		"_csrf":       GetUserCSRFToken(t, session),
		"team_name":   "team1",
		"repo_access": "specific",
		"permission":  "read",
		"unit_8":      "0",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	// user4 can no longer access the org's project
	session = loginUser(t, "user4")
	req = NewRequest(t, "GET", "/org3/-/projects")
	session.MakeRequest(t, req, http.StatusNotFound)
}
