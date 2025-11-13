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

package settings

import (
	"net/http"

	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/services/context"
)

// GetGeneralUISettings returns instance's global settings for ui
func GetGeneralUISettings(ctx *context.APIContext) {
	// swagger:operation GET /settings/ui settings getGeneralUISettings
	// ---
	// summary: Get instance's global settings for ui
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/GeneralUISettings"
	ctx.JSON(http.StatusOK, api.GeneralUISettings{
		DefaultTheme:     setting.UI.DefaultTheme,
		AllowedReactions: setting.UI.Reactions,
		CustomEmojis:     setting.UI.CustomEmojis,
	})
}

// GetGeneralAPISettings returns instance's global settings for api
func GetGeneralAPISettings(ctx *context.APIContext) {
	// swagger:operation GET /settings/api settings getGeneralAPISettings
	// ---
	// summary: Get instance's global settings for api
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/GeneralAPISettings"
	ctx.JSON(http.StatusOK, api.GeneralAPISettings{
		MaxResponseItems:       setting.API.MaxResponseItems,
		DefaultPagingNum:       setting.API.DefaultPagingNum,
		DefaultGitTreesPerPage: setting.API.DefaultGitTreesPerPage,
		DefaultMaxBlobSize:     setting.API.DefaultMaxBlobSize,
		DefaultMaxResponseSize: setting.API.DefaultMaxResponseSize,
	})
}

// GetGeneralRepoSettings returns instance's global settings for repositories
func GetGeneralRepoSettings(ctx *context.APIContext) {
	// swagger:operation GET /settings/repository settings getGeneralRepositorySettings
	// ---
	// summary: Get instance's global settings for repositories
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/GeneralRepoSettings"
	ctx.JSON(http.StatusOK, api.GeneralRepoSettings{
		MirrorsDisabled:      !setting.Mirror.Enabled,
		HTTPGitDisabled:      setting.Repository.DisableHTTPGit,
		MigrationsDisabled:   setting.Repository.DisableMigrations,
		StarsDisabled:        setting.Repository.DisableStars,
		TimeTrackingDisabled: !setting.Service.EnableTimetracking,
		LFSDisabled:          !setting.LFS.StartServer,
	})
}

// GetGeneralAttachmentSettings returns instance's global settings for Attachment
func GetGeneralAttachmentSettings(ctx *context.APIContext) {
	// swagger:operation GET /settings/attachment settings getGeneralAttachmentSettings
	// ---
	// summary: Get instance's global settings for Attachment
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/GeneralAttachmentSettings"
	ctx.JSON(http.StatusOK, api.GeneralAttachmentSettings{
		Enabled:      setting.Attachment.Enabled,
		AllowedTypes: setting.Attachment.AllowedTypes,
		MaxFiles:     setting.Attachment.MaxFiles,
		MaxSize:      setting.Attachment.MaxSize,
	})
}
