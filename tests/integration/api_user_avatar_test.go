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
	"os"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIUpdateUserAvatar(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)

	// Test what happens if you use a valid image
	avatar, err := os.ReadFile("tests/integration/avatar.png")
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "Unable to open avatar.png")
	}

	// Test what happens if you don't have a valid Base64 string
	opts := api.UpdateUserAvatarOption{
		Image: base64.StdEncoding.EncodeToString(avatar),
	}

	req := NewRequestWithJSON(t, "POST", "/api/v1/user/avatar", &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)

	opts = api.UpdateUserAvatarOption{
		Image: "Invalid",
	}

	req = NewRequestWithJSON(t, "POST", "/api/v1/user/avatar", &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusBadRequest)

	// Test what happens if you use a file that is not an image
	text, err := os.ReadFile("tests/integration/README.md")
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "Unable to open README.md")
	}

	opts = api.UpdateUserAvatarOption{
		Image: base64.StdEncoding.EncodeToString(text),
	}

	req = NewRequestWithJSON(t, "POST", "/api/v1/user/avatar", &opts).
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusInternalServerError)
}

func TestAPIDeleteUserAvatar(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)

	req := NewRequest(t, "DELETE", "/api/v1/user/avatar").
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)
}
