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

func getApplyDiffPatchFileOptions() *api.ApplyDiffPatchFileOptions {
	return &api.ApplyDiffPatchFileOptions{
		FileOptions: api.FileOptions{
			BranchName: "master",
		},
		Content: `diff --git a/patch-file-1.txt b/patch-file-1.txt
new file mode 100644
index 0000000000..aaaaaaaaaa
--- /dev/null
+++ b/patch-file-1.txt
@@ -0,0 +1 @@
+File 1
`,
	}
}

func TestAPIApplyDiffPatchFileOptions(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})         // owner of the repo1 & repo16
		org3 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3})          // owner of the repo3, is an org
		user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})         // owner of neither repos
		repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})   // public repo
		repo3 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})   // public repo
		repo16 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 16}) // private repo

		session2 := loginUser(t, user2.Name)
		token2 := getTokenForLoggedInUser(t, session2, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
		session4 := loginUser(t, user4.Name)
		token4 := getTokenForLoggedInUser(t, session4, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)

		req := NewRequestWithJSON(t, "POST", "/api/v1/repos/user2/repo1/diffpatch", getApplyDiffPatchFileOptions()).AddTokenAuth(token2)
		resp := MakeRequest(t, req, http.StatusCreated)
		var fileResponse api.FileResponse
		DecodeJSON(t, resp, &fileResponse)
		assert.Nil(t, fileResponse.Content)
		assert.NotEmpty(t, fileResponse.Commit.HTMLURL)
		req = NewRequest(t, "GET", "/api/v1/repos/user2/repo1/raw/patch-file-1.txt")
		resp = MakeRequest(t, req, http.StatusOK)
		assert.Equal(t, "File 1\n", resp.Body.String())

		// Test creating a file in repo1 by user4 who does not have write access
		req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/diffpatch", user2.Name, repo16.Name), getApplyDiffPatchFileOptions()).
			AddTokenAuth(token4)
		MakeRequest(t, req, http.StatusNotFound)

		// Tests a repo with no token given so will fail
		req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/diffpatch", user2.Name, repo16.Name), getApplyDiffPatchFileOptions())
		MakeRequest(t, req, http.StatusNotFound)

		// Test using access token for a private repo that the user of the token owns
		req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/diffpatch", user2.Name, repo16.Name), getApplyDiffPatchFileOptions()).
			AddTokenAuth(token2)
		MakeRequest(t, req, http.StatusCreated)

		// Test using org repo "org3/repo3" where user2 is a collaborator
		req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/diffpatch", org3.Name, repo3.Name), getApplyDiffPatchFileOptions()).
			AddTokenAuth(token2)
		MakeRequest(t, req, http.StatusCreated)

		// Test using org repo "org3/repo3" with no user token
		req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/diffpatch", org3.Name, repo3.Name), getApplyDiffPatchFileOptions())
		MakeRequest(t, req, http.StatusNotFound)

		// Test using repo "user2/repo1" where user4 is a NOT collaborator
		req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/diffpatch", user2.Name, repo1.Name), getApplyDiffPatchFileOptions()).
			AddTokenAuth(token4)
		MakeRequest(t, req, http.StatusForbidden)
	})
}
