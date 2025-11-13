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

package user

import (
	"io"
	"strings"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/storage"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAvatarLink(t *testing.T) {
	defer test.MockVariableValue(&setting.AppURL, "https://localhost/")()
	defer test.MockVariableValue(&setting.AppSubURL, "")()

	u := &User{ID: 1, Avatar: "avatar.png"}
	link := u.AvatarLink(t.Context())
	assert.Equal(t, "https://localhost/avatars/avatar.png", link)

	setting.AppURL = "https://localhost/sub-path/"
	setting.AppSubURL = "/sub-path"
	link = u.AvatarLink(t.Context())
	assert.Equal(t, "https://localhost/sub-path/avatars/avatar.png", link)
}

func TestUserAvatarGenerate(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	var err error
	tmpDir := t.TempDir()
	storage.Avatars, err = storage.NewLocalStorage(t.Context(), &setting.Storage{Path: tmpDir})
	require.NoError(t, err)

	u := unittest.AssertExistsAndLoadBean(t, &User{ID: 2})

	// there was no avatar, generate a new one
	assert.Empty(t, u.Avatar)
	err = GenerateRandomAvatar(t.Context(), u)
	require.NoError(t, err)
	assert.NotEmpty(t, u.Avatar)

	// make sure the generated one exists
	oldAvatarPath := u.CustomAvatarRelativePath()
	_, err = storage.Avatars.Stat(u.CustomAvatarRelativePath())
	require.NoError(t, err)
	// and try to change its content
	_, err = storage.Avatars.Save(u.CustomAvatarRelativePath(), strings.NewReader("abcd"), 4)
	require.NoError(t, err)

	// try to generate again
	err = GenerateRandomAvatar(t.Context(), u)
	require.NoError(t, err)
	assert.Equal(t, oldAvatarPath, u.CustomAvatarRelativePath())
	f, err := storage.Avatars.Open(u.CustomAvatarRelativePath())
	require.NoError(t, err)
	defer f.Close()
	content, _ := io.ReadAll(f)
	assert.Equal(t, "abcd", string(content))
}
