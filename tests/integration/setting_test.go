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

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestSettingShowUserEmailExplore(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	showUserEmail := setting.UI.ShowUserEmail
	setting.UI.ShowUserEmail = true

	session := loginUser(t, "user2")
	req := NewRequest(t, "GET", "/explore/users?sort=alphabetically")
	resp := session.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, resp.Body)
	assert.Contains(t,
		htmlDoc.doc.Find(".explore.users").Text(),
		"user34@example.com",
	)

	setting.UI.ShowUserEmail = false

	req = NewRequest(t, "GET", "/explore/users?sort=alphabetically")
	resp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, resp.Body)
	assert.NotContains(t,
		htmlDoc.doc.Find(".explore.users").Text(),
		"user34@example.com",
	)

	setting.UI.ShowUserEmail = showUserEmail
}

func TestSettingShowUserEmailProfile(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	showUserEmail := setting.UI.ShowUserEmail

	// user1: keep_email_private = false, user2: keep_email_private = true

	setting.UI.ShowUserEmail = true

	// user1 can see own visible email
	session := loginUser(t, "user1")
	req := NewRequest(t, "GET", "/user1")
	resp := session.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, resp.Body)
	assert.Contains(t, htmlDoc.doc.Find(".user.profile").Text(), "user1@example.com")

	// user1 can not see user2's hidden email
	req = NewRequest(t, "GET", "/user2")
	resp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, resp.Body)
	// Should only contain if the user visits their own profile page
	assert.NotContains(t, htmlDoc.doc.Find(".user.profile").Text(), "user2@example.com")

	// user2 can see user1's visible email
	session = loginUser(t, "user2")
	req = NewRequest(t, "GET", "/user1")
	resp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, resp.Body)
	assert.Contains(t, htmlDoc.doc.Find(".user.profile").Text(), "user1@example.com")

	// user2 can see own hidden email
	session = loginUser(t, "user2")
	req = NewRequest(t, "GET", "/user2")
	resp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, resp.Body)
	assert.Contains(t, htmlDoc.doc.Find(".user.profile").Text(), "user2@example.com")

	setting.UI.ShowUserEmail = false

	// user1 can see own (now hidden) email
	session = loginUser(t, "user1")
	req = NewRequest(t, "GET", "/user1")
	resp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, resp.Body)
	assert.Contains(t, htmlDoc.doc.Find(".user.profile").Text(), "user1@example.com")

	setting.UI.ShowUserEmail = showUserEmail
}

func TestSettingLandingPage(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.LandingPageURL)()

	setting.LandingPageURL = setting.LandingPageHome
	req := NewRequest(t, "GET", "/")
	MakeRequest(t, req, http.StatusOK)

	setting.LandingPageURL = setting.LandingPageExplore
	req = NewRequest(t, "GET", "/")
	resp := MakeRequest(t, req, http.StatusSeeOther)
	assert.Equal(t, "/explore", resp.Header().Get("Location"))

	setting.LandingPageURL = setting.LandingPageOrganizations
	req = NewRequest(t, "GET", "/")
	resp = MakeRequest(t, req, http.StatusSeeOther)
	assert.Equal(t, "/explore/organizations", resp.Header().Get("Location"))

	setting.LandingPageURL = setting.LandingPageLogin
	req = NewRequest(t, "GET", "/")
	resp = MakeRequest(t, req, http.StatusSeeOther)
	assert.Equal(t, "/user/login", resp.Header().Get("Location"))
}
