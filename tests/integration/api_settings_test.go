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
	"net/http"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIExposedSettings(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	ui := new(api.GeneralUISettings)
	req := NewRequest(t, "GET", "/api/v1/settings/ui")
	resp := MakeRequest(t, req, http.StatusOK)

	DecodeJSON(t, resp, &ui)
	assert.Len(t, ui.AllowedReactions, len(setting.UI.Reactions))
	assert.ElementsMatch(t, setting.UI.Reactions, ui.AllowedReactions)

	apiSettings := new(api.GeneralAPISettings)
	req = NewRequest(t, "GET", "/api/v1/settings/api")
	resp = MakeRequest(t, req, http.StatusOK)

	DecodeJSON(t, resp, &apiSettings)
	assert.Equal(t, &api.GeneralAPISettings{
		MaxResponseItems:       setting.API.MaxResponseItems,
		DefaultPagingNum:       setting.API.DefaultPagingNum,
		DefaultGitTreesPerPage: setting.API.DefaultGitTreesPerPage,
		DefaultMaxBlobSize:     setting.API.DefaultMaxBlobSize,
		DefaultMaxResponseSize: setting.API.DefaultMaxResponseSize,
	}, apiSettings)

	repo := new(api.GeneralRepoSettings)
	req = NewRequest(t, "GET", "/api/v1/settings/repository")
	resp = MakeRequest(t, req, http.StatusOK)

	DecodeJSON(t, resp, &repo)
	assert.Equal(t, &api.GeneralRepoSettings{
		MirrorsDisabled:      !setting.Mirror.Enabled,
		HTTPGitDisabled:      setting.Repository.DisableHTTPGit,
		MigrationsDisabled:   setting.Repository.DisableMigrations,
		TimeTrackingDisabled: false,
		LFSDisabled:          !setting.LFS.StartServer,
	}, repo)

	attachment := new(api.GeneralAttachmentSettings)
	req = NewRequest(t, "GET", "/api/v1/settings/attachment")
	resp = MakeRequest(t, req, http.StatusOK)

	DecodeJSON(t, resp, &attachment)
	assert.Equal(t, &api.GeneralAttachmentSettings{
		Enabled:      setting.Attachment.Enabled,
		AllowedTypes: setting.Attachment.AllowedTypes,
		MaxFiles:     setting.Attachment.MaxFiles,
		MaxSize:      setting.Attachment.MaxSize,
	}, attachment)
}
