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

package avatars_test

import (
	"testing"

	avatars_model "github.com/kumose/kmup/models/avatars"
	system_model "github.com/kumose/kmup/models/system"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/setting/config"

	"github.com/stretchr/testify/assert"
)

const gravatarSource = "https://secure.gravatar.com/avatar/"

func disableGravatar(t *testing.T) {
	err := system_model.SetSettings(t.Context(), map[string]string{setting.Config().Picture.EnableFederatedAvatar.DynKey(): "false"})
	assert.NoError(t, err)
	err = system_model.SetSettings(t.Context(), map[string]string{setting.Config().Picture.DisableGravatar.DynKey(): "true"})
	assert.NoError(t, err)
}

func enableGravatar(t *testing.T) {
	err := system_model.SetSettings(t.Context(), map[string]string{setting.Config().Picture.DisableGravatar.DynKey(): "false"})
	assert.NoError(t, err)
	setting.GravatarSource = gravatarSource
}

func TestHashEmail(t *testing.T) {
	assert.Equal(t,
		"d41d8cd98f00b204e9800998ecf8427e",
		avatars_model.HashEmail(""),
	)
	assert.Equal(t,
		"353cbad9b58e69c96154ad99f92bedc7",
		avatars_model.HashEmail("kmup@example.com"),
	)
}

func TestSizedAvatarLink(t *testing.T) {
	setting.AppSubURL = "/testsuburl"

	disableGravatar(t)
	config.GetDynGetter().InvalidateCache()
	assert.Equal(t, "/testsuburl/assets/img/avatar_default.png",
		avatars_model.GenerateEmailAvatarFastLink(t.Context(), "kmup@example.com", 100))

	enableGravatar(t)
	config.GetDynGetter().InvalidateCache()
	assert.Equal(t,
		"https://secure.gravatar.com/avatar/353cbad9b58e69c96154ad99f92bedc7?d=identicon&s=100",
		avatars_model.GenerateEmailAvatarFastLink(t.Context(), "kmup@example.com", 100),
	)
}
