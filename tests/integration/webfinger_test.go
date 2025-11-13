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
	"net/url"
	"strconv"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestWebfinger(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.Federation.Enabled, true)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	appURL, _ := url.Parse(setting.AppURL)

	type webfingerLink struct {
		Rel        string            `json:"rel,omitempty"`
		Type       string            `json:"type,omitempty"`
		Href       string            `json:"href,omitempty"`
		Titles     map[string]string `json:"titles,omitempty"`
		Properties map[string]any    `json:"properties,omitempty"`
	}

	type webfingerJRD struct {
		Subject    string           `json:"subject,omitempty"`
		Aliases    []string         `json:"aliases,omitempty"`
		Properties map[string]any   `json:"properties,omitempty"`
		Links      []*webfingerLink `json:"links,omitempty"`
	}

	session := loginUser(t, "user1")

	req := NewRequest(t, "GET", fmt.Sprintf("/.well-known/webfinger?resource=acct:%s@%s", user.LowerName, appURL.Host))
	resp := MakeRequest(t, req, http.StatusOK)

	var jrd webfingerJRD
	DecodeJSON(t, resp, &jrd)
	assert.Equal(t, "acct:user2@"+appURL.Host, jrd.Subject)
	assert.ElementsMatch(t, []string{user.HTMLURL(t.Context()), appURL.String() + "api/v1/activitypub/user-id/" + strconv.FormatInt(user.ID, 10)}, jrd.Aliases)

	req = NewRequest(t, "GET", fmt.Sprintf("/.well-known/webfinger?resource=acct:%s@%s", user.LowerName, "unknown.host"))
	MakeRequest(t, req, http.StatusBadRequest)

	req = NewRequest(t, "GET", fmt.Sprintf("/.well-known/webfinger?resource=acct:%s@%s", "user31", appURL.Host))
	MakeRequest(t, req, http.StatusNotFound)

	req = NewRequest(t, "GET", fmt.Sprintf("/.well-known/webfinger?resource=acct:%s@%s", "user31", appURL.Host))
	session.MakeRequest(t, req, http.StatusOK)

	req = NewRequest(t, "GET", "/.well-known/webfinger?resource=mailto:"+user.Email)
	MakeRequest(t, req, http.StatusNotFound)
}
