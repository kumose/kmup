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
	"bytes"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/avatar"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestUserAvatar(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2}) // owner of the repo3, is an org

	seed := user2.Email
	if len(seed) == 0 {
		seed = user2.Name
	}

	img, err := avatar.RandomImage([]byte(seed))
	if err != nil {
		assert.NoError(t, err)
		return
	}

	session := loginUser(t, "user2")
	csrf := GetUserCSRFToken(t, session)

	imgData := &bytes.Buffer{}

	body := &bytes.Buffer{}

	// Setup multi-part
	writer := multipart.NewWriter(body)
	writer.WriteField("source", "local")
	part, err := writer.CreateFormFile("avatar", "avatar-for-testuseravatar.png")
	if err != nil {
		assert.NoError(t, err)
		return
	}

	if err := png.Encode(imgData, img); err != nil {
		assert.NoError(t, err)
		return
	}

	if _, err := io.Copy(part, imgData); err != nil {
		assert.NoError(t, err)
		return
	}

	if err := writer.Close(); err != nil {
		assert.NoError(t, err)
		return
	}

	req := NewRequestWithBody(t, "POST", "/user/settings/avatar", body)
	req.Header.Add("X-Csrf-Token", csrf)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	session.MakeRequest(t, req, http.StatusSeeOther)

	user2 = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2}) // owner of the repo3, is an org

	req = NewRequest(t, "GET", user2.AvatarLinkWithSize(t.Context(), 0))
	_ = session.MakeRequest(t, req, http.StatusOK)

	testGetAvatarRedirect(t, user2)

	// Can't test if the response matches because the image is re-generated on upload but checking that this at least doesn't give a 404 should be enough.
}

func testGetAvatarRedirect(t *testing.T, user *user_model.User) {
	t.Run("getAvatarRedirect_"+user.Name, func(t *testing.T) {
		req := NewRequestf(t, "GET", "/%s.png", user.Name)
		resp := MakeRequest(t, req, http.StatusSeeOther)
		assert.Equal(t, "/avatars/"+user.Avatar, resp.Header().Get("location"))
	})
}
