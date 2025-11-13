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

	actions_model "github.com/kumose/kmup/models/actions"
	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestActionsCollaborativeOwner(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		// user2 is the owner of "reusable_workflow" repo
		user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		user2Session := loginUser(t, user2.Name)
		user2Token := getTokenForLoggedInUser(t, user2Session, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
		repo := createActionsTestRepo(t, user2Token, "reusable_workflow", true)

		// a private repo(id=6) of user10 will try to clone "reusable_workflow" repo
		user10 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 10})
		// task id is 55 and its repo_id=6
		task := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionTask{ID: 55, RepoID: 6})
		taskToken := "674f727a81ed2f195bccab036cccf86a182199eb"
		tokenHash := auth_model.HashToken(taskToken, task.TokenSalt)
		assert.Equal(t, task.TokenHash, tokenHash)

		dstPath := t.TempDir()
		u.Path = fmt.Sprintf("%s/%s.git", repo.Owner.UserName, repo.Name)
		u.User = url.UserPassword("kmup-actions", taskToken)

		// the git clone will fail
		doGitCloneFail(u)(t)

		// add user10 to the list of collaborative owners
		req := NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/settings/actions/general/collaborative_owner/add", repo.Owner.UserName, repo.Name), map[string]string{
			"_csrf":               GetUserCSRFToken(t, user2Session),
			"collaborative_owner": user10.Name,
		})
		user2Session.MakeRequest(t, req, http.StatusOK)

		// the git clone will be successful
		doGitClone(dstPath, u)(t)

		// remove user10 from the list of collaborative owners
		req = NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/settings/actions/general/collaborative_owner/delete?id=%d", repo.Owner.UserName, repo.Name, user10.ID), map[string]string{
			"_csrf": GetUserCSRFToken(t, user2Session),
		})
		user2Session.MakeRequest(t, req, http.StatusOK)

		// the git clone will fail
		doGitCloneFail(u)(t)
	})
}
