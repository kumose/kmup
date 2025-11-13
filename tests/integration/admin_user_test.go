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
	"strconv"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAdminViewUsers(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user1")
	req := NewRequest(t, "GET", "/-/admin/users")
	session.MakeRequest(t, req, http.StatusOK)

	session = loginUser(t, "user2")
	req = NewRequest(t, "GET", "/-/admin/users")
	session.MakeRequest(t, req, http.StatusForbidden)
}

func TestAdminViewUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user1")
	req := NewRequest(t, "GET", "/-/admin/users/1")
	session.MakeRequest(t, req, http.StatusOK)

	session = loginUser(t, "user2")
	req = NewRequest(t, "GET", "/-/admin/users/1")
	session.MakeRequest(t, req, http.StatusForbidden)
}

func TestAdminEditUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	testSuccessfullEdit(t, user_model.User{ID: 2, Name: "newusername", LoginName: "otherlogin", Email: "new@e-mail.kmup"})
}

func testSuccessfullEdit(t *testing.T, formData user_model.User) {
	makeRequest(t, formData, http.StatusSeeOther)
}

func makeRequest(t *testing.T, formData user_model.User, headerCode int) {
	session := loginUser(t, "user1")
	csrf := GetUserCSRFToken(t, session)
	req := NewRequestWithValues(t, "POST", "/-/admin/users/"+strconv.Itoa(int(formData.ID))+"/edit", map[string]string{
		"_csrf":      csrf,
		"user_name":  formData.Name,
		"login_name": formData.LoginName,
		"login_type": "0-0",
		"email":      formData.Email,
	})

	session.MakeRequest(t, req, headerCode)
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: formData.ID})
	assert.Equal(t, formData.Name, user.Name)
	assert.Equal(t, formData.LoginName, user.LoginName)
	assert.Equal(t, formData.Email, user.Email)
}

func TestAdminDeleteUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user1")

	usersToDelete := []struct {
		userID int64
		purge  bool
	}{
		{
			userID: 2,
			purge:  true,
		},
		{
			userID: 8,
		},
	}

	for _, entry := range usersToDelete {
		t.Run(fmt.Sprintf("DeleteUser%d", entry.userID), func(t *testing.T) {
			user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: entry.userID})
			assert.NotNil(t, user)

			var query string
			if entry.purge {
				query = "?purge=true"
			}

			csrf := GetUserCSRFToken(t, session)
			req := NewRequestWithValues(t, "POST", fmt.Sprintf("/-/admin/users/%d/delete%s", entry.userID, query), map[string]string{
				"_csrf": csrf,
			})
			session.MakeRequest(t, req, http.StatusSeeOther)

			assertUserDeleted(t, entry.userID)
			unittest.CheckConsistencyFor(t, &user_model.User{})
		})
	}
}
