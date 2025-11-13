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
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	actions_model "github.com/kumose/kmup/models/actions"
	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionsJobTokenAccess(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		t.Run("Write Access", testActionsJobTokenAccess(u, false))
		t.Run("Read Access", testActionsJobTokenAccess(u, true))
	})
}

func testActionsJobTokenAccess(u *url.URL, isFork bool) func(t *testing.T) {
	return func(t *testing.T) {
		task := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionTask{ID: 47})
		require.NoError(t, task.GenerateToken())
		task.Status = actions_model.StatusRunning
		task.IsForkPullRequest = isFork
		err := actions_model.UpdateTask(t.Context(), task, "token_hash", "token_salt", "token_last_eight", "status", "is_fork_pull_request")
		require.NoError(t, err)
		session := emptyTestSession(t)
		context := APITestContext{
			Session:  session,
			Token:    task.Token,
			Username: "user5",
			Reponame: "repo4",
		}
		dstPath := t.TempDir()

		u.Path = context.GitPath()
		u.User = url.UserPassword("kmup-actions", task.Token)

		t.Run("Git Clone", doGitClone(dstPath, u))

		t.Run("API Get Repository", doAPIGetRepository(context, func(t *testing.T, r structs.Repository) {
			require.Equal(t, "repo4", r.Name)
			require.Equal(t, "user5", r.Owner.UserName)
		}))

		context.ExpectedCode = util.Iif(isFork, http.StatusForbidden, http.StatusCreated)
		t.Run("API Create File", doAPICreateFile(context, "test.txt", &structs.CreateFileOptions{
			FileOptions: structs.FileOptions{
				NewBranchName: "new-branch",
				Message:       "Create File",
			},
			ContentBase64: base64.StdEncoding.EncodeToString([]byte(`This is a test file created using job token.`)),
		}))

		context.ExpectedCode = http.StatusForbidden
		t.Run("Fail to Create Repository", doAPICreateRepository(context, true))

		context.ExpectedCode = http.StatusForbidden
		t.Run("Fail to Delete Repository", doAPIDeleteRepository(context))

		t.Run("Fail to Create Organization", doAPICreateOrganization(context, &structs.CreateOrgOption{
			UserName: "actions",
			FullName: "Kmup Actions",
		}))
	}
}

func TestActionsJobTokenAccessLFS(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		httpContext := NewAPITestContext(t, "user2", "repo-lfs-test", auth_model.AccessTokenScopeWriteUser, auth_model.AccessTokenScopeWriteRepository)
		t.Run("Create Repository", doAPICreateRepository(httpContext, false, func(t *testing.T, repository structs.Repository) {
			task := &actions_model.ActionTask{}
			require.NoError(t, task.GenerateToken())
			task.Status = actions_model.StatusRunning
			task.IsForkPullRequest = false
			task.RepoID = repository.ID
			err := db.Insert(t.Context(), task)
			require.NoError(t, err)
			session := emptyTestSession(t)
			httpContext := APITestContext{
				Session:  session,
				Token:    task.Token,
				Username: "user2",
				Reponame: "repo-lfs-test",
			}

			u.Path = httpContext.GitPath()
			dstPath := t.TempDir()

			u.Path = httpContext.GitPath()
			u.User = url.UserPassword("kmup-actions", task.Token)

			t.Run("Clone", doGitClone(dstPath, u))

			dstPath2 := t.TempDir()

			t.Run("Partial Clone", doPartialGitClone(dstPath2, u))

			lfs := lfsCommitAndPushTest(t, dstPath, testFileSizeSmall)[0]

			reqLFS := NewRequest(t, "GET", "/api/v1/repos/user2/repo-lfs-test/media/"+lfs).AddTokenAuth(task.Token)
			respLFS := MakeRequestNilResponseRecorder(t, reqLFS, http.StatusOK)
			assert.Equal(t, testFileSizeSmall, respLFS.Length)
		}))
	})
}
