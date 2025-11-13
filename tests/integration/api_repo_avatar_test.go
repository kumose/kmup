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
	"fmt"
	"net/http"
	"os"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIUpdateRepoAvatar(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	token := getUserToken(t, user2.LowerName, auth_model.AccessTokenScopeWriteRepository)

	// Test what happens if you use a valid image
	avatar, err := os.ReadFile("tests/integration/avatar.png")
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "Unable to open avatar.png")
	}

	opts := api.UpdateRepoAvatarOption{
		Image: base64.StdEncoding.EncodeToString(avatar),
	}

	req := NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/avatar", repo.OwnerName, repo.Name), &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	// Test what happens if you don't have a valid Base64 string
	opts = api.UpdateRepoAvatarOption{
		Image: "Invalid",
	}

	req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/avatar", repo.OwnerName, repo.Name), &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusBadRequest)

	// Test what happens if you use a file that is not an image
	text, err := os.ReadFile("tests/integration/README.md")
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "Unable to open README.md")
	}

	opts = api.UpdateRepoAvatarOption{
		Image: base64.StdEncoding.EncodeToString(text),
	}

	req = NewRequestWithJSON(t, "POST", fmt.Sprintf("/api/v1/repos/%s/%s/avatar", repo.OwnerName, repo.Name), &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusInternalServerError)
}

func TestAPIDeleteRepoAvatar(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	token := getUserToken(t, user2.LowerName, auth_model.AccessTokenScopeWriteRepository)

	req := NewRequest(t, "DELETE", fmt.Sprintf("/api/v1/repos/%s/%s/avatar", repo.OwnerName, repo.Name)).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)
}
