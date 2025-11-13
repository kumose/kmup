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
	"testing"

	git_model "github.com/kumose/kmup/models/git"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestRenameBranch(t *testing.T) {
	onKmupRun(t, testRenameBranch)
}

func testRenameBranch(t *testing.T, u *url.URL) {
	defer tests.PrepareTestEnv(t)()

	unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: 1, Name: "master"})

	// get branch setting page
	session := loginUser(t, "user2")
	req := NewRequest(t, "GET", "/user2/repo1/branches")
	resp := session.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, resp.Body)

	req = NewRequestWithValues(t, "POST", "/user2/repo1/branches/rename", map[string]string{
		"_csrf": htmlDoc.GetCSRF(),
		"from":  "master",
		"to":    "master",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	// check new branch link
	req = NewRequestWithValues(t, "GET", "/user2/repo1/src/branch/master/README.md", nil)
	session.MakeRequest(t, req, http.StatusOK)

	// check old branch link
	req = NewRequestWithValues(t, "GET", "/user2/repo1/src/branch/master/README.md", nil)
	resp = session.MakeRequest(t, req, http.StatusSeeOther)
	location := resp.Header().Get("Location")
	assert.Equal(t, "/user2/repo1/src/branch/master/README.md", location)

	// check db
	repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.Equal(t, "master", repo1.DefaultBranch)

	// create branch1
	csrf := GetUserCSRFToken(t, session)

	req = NewRequestWithValues(t, "POST", "/user2/repo1/branches/_new/branch/master", map[string]string{
		"_csrf":           csrf,
		"new_branch_name": "branch1",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	branch1 := unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch1"})
	assert.Equal(t, "branch1", branch1.Name)

	// create branch2
	req = NewRequestWithValues(t, "POST", "/user2/repo1/branches/_new/branch/master", map[string]string{
		"_csrf":           csrf,
		"new_branch_name": "branch2",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	branch2 := unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch2"})
	assert.Equal(t, "branch2", branch2.Name)

	// rename branch2 to branch1
	req = NewRequestWithValues(t, "POST", "/user2/repo1/branches/rename", map[string]string{
		"_csrf": htmlDoc.GetCSRF(),
		"from":  "branch2",
		"to":    "branch1",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)
	flashMsg := session.GetCookieFlashMessage()
	assert.NotEmpty(t, flashMsg.ErrorMsg)

	branch2 = unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch2"})
	assert.Equal(t, "branch2", branch2.Name)
	branch1 = unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch1"})
	assert.Equal(t, "branch1", branch1.Name)

	// delete branch1
	req = NewRequestWithValues(t, "POST", "/user2/repo1/branches/delete", map[string]string{
		"_csrf": htmlDoc.GetCSRF(),
		"name":  "branch1",
	})
	session.MakeRequest(t, req, http.StatusOK)
	branch2 = unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch2"})
	assert.Equal(t, "branch2", branch2.Name)
	branch1 = unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch1"})
	assert.True(t, branch1.IsDeleted) // virtual deletion

	// rename branch2 to branch1 again
	req = NewRequestWithValues(t, "POST", "/user2/repo1/branches/rename", map[string]string{
		"_csrf": htmlDoc.GetCSRF(),
		"from":  "branch2",
		"to":    "branch1",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	flashMsg = session.GetCookieFlashMessage()
	assert.NotEmpty(t, flashMsg.SuccessMsg)

	unittest.AssertNotExistsBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch2"})
	branch1 = unittest.AssertExistsAndLoadBean(t, &git_model.Branch{RepoID: repo1.ID, Name: "branch1"})
	assert.Equal(t, "branch1", branch1.Name)
}
