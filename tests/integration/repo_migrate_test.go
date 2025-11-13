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
	"testing"

	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func testRepoMigrate(t testing.TB, session *TestSession, cloneAddr, repoName string) *httptest.ResponseRecorder {
	req := NewRequest(t, "GET", fmt.Sprintf("/repo/migrate?service_type=%d", structs.PlainGitService)) // render plain git migration page
	resp := session.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, resp.Body)

	link, exists := htmlDoc.doc.Find("form.ui.form").Attr("action")
	assert.True(t, exists, "The template has changed")

	uid, exists := htmlDoc.doc.Find("#uid").Attr("value")
	assert.True(t, exists, "The template has changed")

	req = NewRequestWithValues(t, "POST", link, map[string]string{
		"_csrf":      htmlDoc.GetCSRF(),
		"clone_addr": cloneAddr,
		"uid":        uid,
		"repo_name":  repoName,
		"service":    fmt.Sprintf("%d", structs.PlainGitService),
	})
	resp = session.MakeRequest(t, req, http.StatusSeeOther)

	return resp
}

func TestRepoMigrate(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	session := loginUser(t, "user2")
	testRepoMigrate(t, session, "https://github.com/go-kmup/test_repo.git", "git")
}
