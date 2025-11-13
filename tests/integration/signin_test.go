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
	"net/url"
	"strings"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/modules/translation"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers"
	"github.com/kumose/kmup/routers/web/auth"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/tests"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLoginFailed(t *testing.T, username, password, message string) {
	session := emptyTestSession(t)
	req := NewRequestWithValues(t, "POST", "/user/login", map[string]string{
		"user_name": username,
		"password":  password,
	})
	resp := session.MakeRequest(t, req, http.StatusOK)

	htmlDoc := NewHTMLParser(t, resp.Body)
	resultMsg := htmlDoc.doc.Find(".ui.message>p").Text()

	assert.Equal(t, message, resultMsg)
}

func TestSignin(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	// add new user with user2's email
	user.Name = "testuser"
	user.LowerName = strings.ToLower(user.Name)
	user.ID = 0
	require.NoError(t, db.Insert(t.Context(), user))

	samples := []struct {
		username string
		password string
		message  string
	}{
		{username: "wrongUsername", password: "wrongPassword", message: translation.NewLocale("en-US").TrString("form.username_password_incorrect")},
		{username: "wrongUsername", password: "password", message: translation.NewLocale("en-US").TrString("form.username_password_incorrect")},
		{username: "user15", password: "wrongPassword", message: translation.NewLocale("en-US").TrString("form.username_password_incorrect")},
		{username: "user1@example.com", password: "wrongPassword", message: translation.NewLocale("en-US").TrString("form.username_password_incorrect")},
	}

	for _, s := range samples {
		testLoginFailed(t, s.username, s.password, s.message)
	}
}

func TestSigninWithRememberMe(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	baseURL, _ := url.Parse(setting.AppURL)

	session := emptyTestSession(t)
	req := NewRequestWithValues(t, "POST", "/user/login", map[string]string{
		"user_name": user.Name,
		"password":  userPassword,
		"remember":  "on",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	c := session.GetRawCookie(setting.CookieRememberName)
	assert.NotNil(t, c)

	session = emptyTestSession(t)

	// Without session the settings page should not be reachable
	req = NewRequest(t, "GET", "/user/settings")
	session.MakeRequest(t, req, http.StatusSeeOther)

	req = NewRequest(t, "GET", "/user/login")
	// Set the remember me cookie for the login GET request
	session.jar.SetCookies(baseURL, []*http.Cookie{c})
	session.MakeRequest(t, req, http.StatusSeeOther)

	// With session the settings page should be reachable
	req = NewRequest(t, "GET", "/user/settings")
	session.MakeRequest(t, req, http.StatusOK)
}

func TestEnablePasswordSignInFormAndEnablePasskeyAuth(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	mockLinkAccount := func(ctx *context.Context) {
		authSource := auth_model.Source{ID: 1}
		gothUser := goth.User{Email: "invalid-email", Name: "."}
		_ = auth.Oauth2SetLinkAccountData(ctx, auth.LinkAccountData{AuthSourceID: authSource.ID, GothUser: gothUser})
	}

	t.Run("EnablePasswordSignInForm=false", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		defer test.MockVariableValue(&setting.Service.EnablePasswordSignInForm, false)()

		req := NewRequest(t, "GET", "/user/login")
		resp := MakeRequest(t, req, http.StatusOK)
		doc := NewHTMLParser(t, resp.Body)
		AssertHTMLElement(t, doc, "form[action='/user/login']", false)

		req = NewRequest(t, "POST", "/user/login")
		MakeRequest(t, req, http.StatusForbidden)

		req = NewRequest(t, "GET", "/user/link_account")
		defer web.RouteMockReset()
		web.RouteMock(web.MockAfterMiddlewares, mockLinkAccount)
		resp = MakeRequest(t, req, http.StatusOK)
		doc = NewHTMLParser(t, resp.Body)
		AssertHTMLElement(t, doc, "form[action='/user/link_account_signin']", false)
	})

	t.Run("EnablePasswordSignInForm=true", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		defer test.MockVariableValue(&setting.Service.EnablePasswordSignInForm, true)()

		req := NewRequest(t, "GET", "/user/login")
		resp := MakeRequest(t, req, http.StatusOK)
		doc := NewHTMLParser(t, resp.Body)
		AssertHTMLElement(t, doc, "form[action='/user/login']", true)

		req = NewRequest(t, "POST", "/user/login")
		MakeRequest(t, req, http.StatusOK)

		req = NewRequest(t, "GET", "/user/link_account")
		defer web.RouteMockReset()
		web.RouteMock(web.MockAfterMiddlewares, mockLinkAccount)
		resp = MakeRequest(t, req, http.StatusOK)
		doc = NewHTMLParser(t, resp.Body)
		AssertHTMLElement(t, doc, "form[action='/user/link_account_signin']", true)
	})

	t.Run("EnablePasskeyAuth=false", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		defer test.MockVariableValue(&setting.Service.EnablePasskeyAuth, false)()

		req := NewRequest(t, "GET", "/user/login")
		resp := MakeRequest(t, req, http.StatusOK)
		doc := NewHTMLParser(t, resp.Body)
		AssertHTMLElement(t, doc, ".signin-passkey", false)
	})

	t.Run("EnablePasskeyAuth=true", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		defer test.MockVariableValue(&setting.Service.EnablePasskeyAuth, true)()

		req := NewRequest(t, "GET", "/user/login")
		resp := MakeRequest(t, req, http.StatusOK)
		doc := NewHTMLParser(t, resp.Body)
		AssertHTMLElement(t, doc, ".signin-passkey", true)
	})
}

func TestRequireSignInView(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	t.Run("NoRequireSignInView", func(t *testing.T) {
		require.False(t, setting.Service.RequireSignInViewStrict)
		require.False(t, setting.Service.BlockAnonymousAccessExpensive)
		req := NewRequest(t, "GET", "/user2/repo1/src/branch/master")
		MakeRequest(t, req, http.StatusOK)
	})
	t.Run("RequireSignInView", func(t *testing.T) {
		defer test.MockVariableValue(&setting.Service.RequireSignInViewStrict, true)()
		defer test.MockVariableValue(&testWebRoutes, routers.NormalRoutes())()
		req := NewRequest(t, "GET", "/user2/repo1/src/branch/master")
		resp := MakeRequest(t, req, http.StatusSeeOther)
		assert.Equal(t, "/user/login", resp.Header().Get("Location"))
	})
	t.Run("BlockAnonymousAccessExpensive", func(t *testing.T) {
		defer test.MockVariableValue(&setting.Service.RequireSignInViewStrict, false)()
		defer test.MockVariableValue(&setting.Service.BlockAnonymousAccessExpensive, true)()
		defer test.MockVariableValue(&testWebRoutes, routers.NormalRoutes())()

		req := NewRequest(t, "GET", "/user2/repo1")
		MakeRequest(t, req, http.StatusOK)

		req = NewRequest(t, "GET", "/user2/repo1/src/branch/master")
		resp := MakeRequest(t, req, http.StatusSeeOther)
		assert.Equal(t, "/user/login", resp.Header().Get("Location"))
	})
}
