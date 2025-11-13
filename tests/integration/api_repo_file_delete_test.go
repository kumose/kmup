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

	auth_model "github.com/kumose/kmup/models/auth"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"

	"github.com/stretchr/testify/assert"
)

func getDeleteFileOptions() *api.DeleteFileOptions {
	return &api.DeleteFileOptions{
		SHA: "103ff9234cefeee5ec5361d22b49fbb04d385885",
		FileOptions: api.FileOptions{
			BranchName:    "master",
			NewBranchName: "master",
			Message:       "Removing the file new/file.txt",
			Author: api.Identity{
				Name:  "John Doe",
				Email: "johndoe@example.com",
			},
			Committer: api.Identity{
				Name:  "Jane Doe",
				Email: "janedoe@example.com",
			},
		},
	}
}

func TestAPIDeleteFile(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})         // owner of the repo1 & repo16
		org3 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3})          // owner of the repo3, is an org
		user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})         // owner of neither repos
		repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})   // public repo
		repo3 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})   // public repo
		repo16 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 16}) // private repo
		fileID := 0

		// Get user2's token
		session := loginUser(t, user2.Name)
		token2 := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
		// Get user4's token
		session = loginUser(t, user4.Name)
		token4 := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

		// Test deleting a file in repo1 which user2 owns, try both with branch and empty branch
		for _, branch := range [...]string{
			"master", // Branch
			"",       // Empty branch
		} {
			fileID++
			treePath := fmt.Sprintf("delete/file%d.txt", fileID)
			createFile(user2, repo1, treePath)
			deleteFileOptions := getDeleteFileOptions()
			deleteFileOptions.BranchName = branch
			req := NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo1.Name, treePath), &deleteFileOptions).
				AddTokenAuth(token2)
			resp := MakeRequest(t, req, http.StatusOK)
			var fileResponse api.FileResponse
			DecodeJSON(t, resp, &fileResponse)
			assert.NotNil(t, fileResponse)
			assert.Nil(t, fileResponse.Content)
		}

		// Test deleting file and making the delete in a new branch
		fileID++
		treePath := fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(user2, repo1, treePath)
		deleteFileOptions := getDeleteFileOptions()
		deleteFileOptions.BranchName = repo1.DefaultBranch
		deleteFileOptions.NewBranchName = "new_branch"
		req := NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo1.Name, treePath), &deleteFileOptions).
			AddTokenAuth(token2)
		resp := MakeRequest(t, req, http.StatusOK)
		var fileResponse api.FileResponse
		DecodeJSON(t, resp, &fileResponse)
		assert.NotNil(t, fileResponse)
		assert.Nil(t, fileResponse.Content)
		assert.Equal(t, deleteFileOptions.Message+"\n", fileResponse.Commit.Message)

		// Test deleting file without a message
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(user2, repo1, treePath)
		deleteFileOptions = getDeleteFileOptions()
		deleteFileOptions.Message = ""
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo1.Name, treePath), &deleteFileOptions).
			AddTokenAuth(token2)
		resp = MakeRequest(t, req, http.StatusOK)
		DecodeJSON(t, resp, &fileResponse)
		expectedMessage := "Delete " + treePath + "\n"
		assert.Equal(t, expectedMessage, fileResponse.Commit.Message)

		// Test deleting a file with the wrong SHA
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(user2, repo1, treePath)
		deleteFileOptions = getDeleteFileOptions()
		deleteFileOptions.SHA = "badsha"
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo1.Name, treePath), &deleteFileOptions).
			AddTokenAuth(token2)
		MakeRequest(t, req, http.StatusUnprocessableEntity)

		// Test creating a file in repo16 by user4 who does not have write access
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(user2, repo16, treePath)
		deleteFileOptions = getDeleteFileOptions()
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo16.Name, treePath), &deleteFileOptions).
			AddTokenAuth(token4)
		MakeRequest(t, req, http.StatusNotFound)

		// Tests a repo with no token given so will fail
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(user2, repo16, treePath)
		deleteFileOptions = getDeleteFileOptions()
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo16.Name, treePath), &deleteFileOptions)
		MakeRequest(t, req, http.StatusNotFound)

		// Test using access token for a private repo that the user of the token owns
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(user2, repo16, treePath)
		deleteFileOptions = getDeleteFileOptions()
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo16.Name, treePath), &deleteFileOptions).
			AddTokenAuth(token2)
		MakeRequest(t, req, http.StatusOK)

		// Test using org repo "org3/repo3" where user2 is a collaborator
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(org3, repo3, treePath)
		deleteFileOptions = getDeleteFileOptions()
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", org3.Name, repo3.Name, treePath), &deleteFileOptions).
			AddTokenAuth(token2)
		MakeRequest(t, req, http.StatusOK)

		// Test using org repo "org3/repo3" with no user token
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(org3, repo3, treePath)
		deleteFileOptions = getDeleteFileOptions()
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", org3.Name, repo3.Name, treePath), &deleteFileOptions)
		MakeRequest(t, req, http.StatusNotFound)

		// Test using repo "user2/repo1" where user4 is a NOT collaborator
		fileID++
		treePath = fmt.Sprintf("delete/file%d.txt", fileID)
		createFile(user2, repo1, treePath)
		deleteFileOptions = getDeleteFileOptions()
		req = NewRequestWithJSON(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", user2.Name, repo1.Name, treePath), &deleteFileOptions).
			AddTokenAuth(token4)
		MakeRequest(t, req, http.StatusForbidden)
	})
}
