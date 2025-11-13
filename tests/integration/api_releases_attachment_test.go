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
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/tests"
)

func TestAPIEditReleaseAttachmentWithUnallowedFile(t *testing.T) {
	// Limit the allowed release types (since by default there is no restriction)
	defer test.MockVariableValue(&setting.Repository.Release.AllowedTypes, ".exe")()
	defer tests.PrepareTestEnv(t)()

	attachment := unittest.AssertExistsAndLoadBean(t, &repo_model.Attachment{ID: 9})
	release := unittest.AssertExistsAndLoadBean(t, &repo_model.Release{ID: attachment.ReleaseID})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: attachment.RepoID})
	repoOwner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, repoOwner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

	filename := "file.bad"
	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/releases/%d/assets/%d", repoOwner.Name, repo.Name, release.ID, attachment.ID)
	req := NewRequestWithValues(t, "PATCH", urlStr, map[string]string{
		"name": filename,
	}).AddTokenAuth(token)

	session.MakeRequest(t, req, http.StatusUnprocessableEntity)
}
