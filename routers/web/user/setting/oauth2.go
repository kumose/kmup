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

package setting

import (
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
)

const (
	tplSettingsOAuthApplicationEdit templates.TplName = "user/settings/applications_oauth2_edit"
)

func newOAuth2CommonHandlers(userID int64) *OAuth2CommonHandlers {
	return &OAuth2CommonHandlers{
		OwnerID:            userID,
		BasePathList:       setting.AppSubURL + "/user/settings/applications",
		BasePathEditPrefix: setting.AppSubURL + "/user/settings/applications/oauth2",
		TplAppEdit:         tplSettingsOAuthApplicationEdit,
	}
}

// OAuthApplicationsPost response for adding a oauth2 application
func OAuthApplicationsPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsApplications"] = true

	oa := newOAuth2CommonHandlers(ctx.Doer.ID)
	oa.AddApp(ctx)
}

// OAuthApplicationsEdit response for editing oauth2 application
func OAuthApplicationsEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsApplications"] = true

	oa := newOAuth2CommonHandlers(ctx.Doer.ID)
	oa.EditSave(ctx)
}

// OAuthApplicationsRegenerateSecret handles the post request for regenerating the secret
func OAuthApplicationsRegenerateSecret(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsApplications"] = true

	oa := newOAuth2CommonHandlers(ctx.Doer.ID)
	oa.RegenerateSecret(ctx)
}

// OAuth2ApplicationShow displays the given application
func OAuth2ApplicationShow(ctx *context.Context) {
	oa := newOAuth2CommonHandlers(ctx.Doer.ID)
	oa.EditShow(ctx)
}

// DeleteOAuth2Application deletes the given oauth2 application
func DeleteOAuth2Application(ctx *context.Context) {
	oa := newOAuth2CommonHandlers(ctx.Doer.ID)
	oa.DeleteApp(ctx)
}

// RevokeOAuth2Grant revokes the grant with the given id
func RevokeOAuth2Grant(ctx *context.Context) {
	oa := newOAuth2CommonHandlers(ctx.Doer.ID)
	oa.RevokeGrant(ctx)
}
