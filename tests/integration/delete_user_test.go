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
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/organization"
	access_model "github.com/kumose/kmup/models/perm/access"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/tests"
)

func assertUserDeleted(t *testing.T, userID int64) {
	unittest.AssertNotExistsBean(t, &user_model.User{ID: userID})
	unittest.AssertNotExistsBean(t, &user_model.Follow{UserID: userID})
	unittest.AssertNotExistsBean(t, &user_model.Follow{FollowID: userID})
	unittest.AssertNotExistsBean(t, &repo_model.Repository{OwnerID: userID})
	unittest.AssertNotExistsBean(t, &access_model.Access{UserID: userID})
	unittest.AssertNotExistsBean(t, &organization.OrgUser{UID: userID})
	unittest.AssertNotExistsBean(t, &issues_model.IssueUser{UID: userID})
	unittest.AssertNotExistsBean(t, &organization.TeamUser{UID: userID})
	unittest.AssertNotExistsBean(t, &repo_model.Star{UID: userID})
}

func TestUserDeleteAccount(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user8")
	csrf := GetUserCSRFToken(t, session)
	urlStr := "/user/settings/account/delete?password=" + userPassword
	req := NewRequestWithValues(t, "POST", urlStr, map[string]string{
		"_csrf": csrf,
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	assertUserDeleted(t, 8)
	unittest.CheckConsistencyFor(t, &user_model.User{})
}

func TestUserDeleteAccountStillOwnRepos(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")
	csrf := GetUserCSRFToken(t, session)
	urlStr := "/user/settings/account/delete?password=" + userPassword
	req := NewRequestWithValues(t, "POST", urlStr, map[string]string{
		"_csrf": csrf,
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	// user should not have been deleted, because the user still owns repos
	unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
}
