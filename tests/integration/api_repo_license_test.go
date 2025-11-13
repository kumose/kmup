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
	"testing"
	"time"

	auth_model "github.com/kumose/kmup/models/auth"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

var testLicenseContent = `
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
`

func TestAPIRepoLicense(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		session := loginUser(t, "user2")

		// Request editor page
		req := NewRequest(t, "GET", "/user2/repo1/_new/master/")
		resp := session.MakeRequest(t, req, http.StatusOK)

		doc := NewHTMLParser(t, resp.Body)
		lastCommit := doc.GetInputValueByName("last_commit")
		assert.NotEmpty(t, lastCommit)

		// Save new file to master branch
		req = NewRequestWithValues(t, "POST", "/user2/repo1/_new/master/", map[string]string{
			"_csrf":         doc.GetCSRF(),
			"last_commit":   lastCommit,
			"tree_path":     "LICENSE",
			"content":       testLicenseContent,
			"commit_choice": "direct",
		})
		resp = session.MakeRequest(t, req, http.StatusOK)
		assert.NotEmpty(t, test.RedirectURL(resp))

		// let kmup update repo license
		time.Sleep(time.Second)
		checkRepoLicense(t, "user2", "repo1", []string{"BSD-2-Clause"})

		// Change default branch
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
		branchName := "DefaultBranch"
		req = NewRequestWithJSON(t, "PATCH", "/api/v1/repos/user2/repo1", api.EditRepoOption{
			DefaultBranch: &branchName,
		}).AddTokenAuth(token)
		session.MakeRequest(t, req, http.StatusOK)

		// let kmup update repo license
		time.Sleep(time.Second)
		checkRepoLicense(t, "user2", "repo1", []string{"MIT"})
	})
}

func checkRepoLicense(t *testing.T, owner, repo string, expected []string) {
	reqURL := fmt.Sprintf("/api/v1/repos/%s/%s/licenses", owner, repo)
	req := NewRequest(t, "GET", reqURL)
	resp := MakeRequest(t, req, http.StatusOK)

	var licenses []string
	DecodeJSON(t, resp, &licenses)

	assert.ElementsMatch(t, expected, licenses, 0)
}
