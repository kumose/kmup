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
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func testRepoGenerate(t *testing.T, session *TestSession, templateID, templateOwnerName, templateRepoName, generateOwnerName, generateRepoName string) *httptest.ResponseRecorder {
	generateOwner := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: generateOwnerName})

	// Step0: check the existence of the generated repo
	req := NewRequestf(t, "GET", "/%s/%s", generateOwnerName, generateRepoName)
	session.MakeRequest(t, req, http.StatusNotFound)

	// Step1: go to the main page of template repo
	req = NewRequestf(t, "GET", "/%s/%s", templateOwnerName, templateRepoName)
	resp := session.MakeRequest(t, req, http.StatusOK)

	// Step2: click the "Use this template" button
	htmlDoc := NewHTMLParser(t, resp.Body)
	link, exists := htmlDoc.doc.Find(`a.ui.button[href^="/repo/create"]`).Attr("href")
	assert.True(t, exists, "The template has changed")
	req = NewRequest(t, "GET", link)
	resp = session.MakeRequest(t, req, http.StatusOK)

	// Step3: fill the form on the "create" page
	htmlDoc = NewHTMLParser(t, resp.Body)
	link, exists = htmlDoc.doc.Find(`form.ui.form[action^="/repo/create"]`).Attr("action")
	assert.True(t, exists, "The template has changed")
	_, exists = htmlDoc.doc.Find(fmt.Sprintf(`#repo_owner_dropdown .item[data-value="%d"]`, generateOwner.ID)).Attr("data-value")
	assert.True(t, exists, "Generate owner '%s' is not present in select box", generateOwnerName)
	req = NewRequestWithValues(t, "POST", link, map[string]string{
		"_csrf":         htmlDoc.GetCSRF(),
		"uid":           strconv.FormatInt(generateOwner.ID, 10),
		"repo_name":     generateRepoName,
		"repo_template": templateID,
		"git_content":   "true",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	// Step4: check the existence of the generated repo
	req = NewRequestf(t, "GET", "/%s/%s", generateOwnerName, generateRepoName)
	session.MakeRequest(t, req, http.StatusOK)

	// Step5: check substituted values in Readme
	req = NewRequestf(t, "GET", "/%s/%s/raw/branch/master/README.md", generateOwnerName, generateRepoName)
	resp = session.MakeRequest(t, req, http.StatusOK)
	body := fmt.Sprintf(`# %s Readme
Owner: %s
Link: /%s/%s
Clone URL: %s%s/%s.git`,
		generateRepoName,
		strings.ToUpper(generateOwnerName),
		generateOwnerName,
		generateRepoName,
		setting.AppURL,
		generateOwnerName,
		generateRepoName)
	assert.Equal(t, body, resp.Body.String())

	// Step6: check substituted values in substituted file path ${REPO_NAME}
	req = NewRequestf(t, "GET", "/%s/%s/raw/branch/master/%s.log", generateOwnerName, generateRepoName, generateRepoName)
	resp = session.MakeRequest(t, req, http.StatusOK)
	assert.Equal(t, generateRepoName, resp.Body.String())

	return resp
}

func TestRepoGenerate(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	session := loginUser(t, "user1")
	testRepoGenerate(t, session, "44", "user27", "template1", "user1", "generated1")
}

func TestRepoGenerateToOrg(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	session := loginUser(t, "user2")
	testRepoGenerate(t, session, "44", "user27", "template1", "user2", "generated2")
}
