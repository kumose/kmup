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

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestXSSUserFullName(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	const fullName = `name & <script class="evil">alert('Oh no!');</script>`

	session := loginUser(t, user.Name)
	req := NewRequestWithValues(t, "POST", "/user/settings", map[string]string{
		"_csrf":     GetUserCSRFToken(t, session),
		"name":      user.Name,
		"full_name": fullName,
		"email":     user.Email,
		"language":  "en-US",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	req = NewRequestf(t, "GET", "/%s", user.Name)
	resp := session.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, resp.Body)
	assert.Equal(t, 0, htmlDoc.doc.Find("script.evil").Length())
	assert.Equal(t, fullName,
		htmlDoc.doc.Find("div.content").Find(".header.text.center").Text(),
	)
}
