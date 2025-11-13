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
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIRepoTags(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

	repoName := "repo1"

	req := NewRequestf(t, "GET", "/api/v1/repos/%s/%s/tags", user.Name, repoName).
		AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusOK)

	var tags []*api.Tag
	DecodeJSON(t, resp, &tags)

	assert.Len(t, tags, 1)
	assert.Equal(t, "v1.1", tags[0].Name)
	assert.Equal(t, "Initial commit", tags[0].Message)
	assert.Equal(t, "65f1bf27bc3bf70f64657658635e66094edbcb4d", tags[0].Commit.SHA)
	assert.Equal(t, setting.AppURL+"api/v1/repos/user2/repo1/git/commits/65f1bf27bc3bf70f64657658635e66094edbcb4d", tags[0].Commit.URL)
	assert.Equal(t, setting.AppURL+"user2/repo1/archive/v1.1.zip", tags[0].ZipballURL)
	assert.Equal(t, setting.AppURL+"user2/repo1/archive/v1.1.tar.gz", tags[0].TarballURL)

	newTag := createNewTagUsingAPI(t, token, user.Name, repoName, "kmup/22", "", "nice!\nand some text")
	resp = MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, &tags)
	assert.Len(t, tags, 2)
	for _, tag := range tags {
		if tag.Name != "v1.1" {
			assert.Equal(t, newTag.Name, tag.Name)
			assert.Equal(t, newTag.Message, tag.Message)
			assert.Equal(t, "nice!\nand some text", tag.Message)
			assert.Equal(t, newTag.Commit.SHA, tag.Commit.SHA)
		}
	}

	// get created tag
	req = NewRequestf(t, "GET", "/api/v1/repos/%s/%s/tags/%s", user.Name, repoName, newTag.Name).
		AddTokenAuth(token)
	resp = MakeRequest(t, req, http.StatusOK)
	var tag *api.Tag
	DecodeJSON(t, resp, &tag)
	assert.Equal(t, newTag, tag)

	// delete tag
	delReq := NewRequestf(t, "DELETE", "/api/v1/repos/%s/%s/tags/%s", user.Name, repoName, newTag.Name).
		AddTokenAuth(token)
	MakeRequest(t, delReq, http.StatusNoContent)

	// check if it's gone
	MakeRequest(t, req, http.StatusNotFound)
}

func createNewTagUsingAPI(t *testing.T, token, ownerName, repoName, name, target, msg string) *api.Tag {
	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/tags", ownerName, repoName)
	req := NewRequestWithJSON(t, "POST", urlStr, &api.CreateTagOption{
		TagName: name,
		Message: msg,
		Target:  target,
	}).AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusCreated)

	var respObj api.Tag
	DecodeJSON(t, resp, &respObj)
	return &respObj
}
