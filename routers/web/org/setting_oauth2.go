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

package org

import (
	"fmt"
	"net/http"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	user_setting "github.com/kumose/kmup/routers/web/user/setting"
	"github.com/kumose/kmup/services/context"
)

const (
	tplSettingsApplications         templates.TplName = "org/settings/applications"
	tplSettingsOAuthApplicationEdit templates.TplName = "org/settings/applications_oauth2_edit"
)

func newOAuth2CommonHandlers(org *context.Organization) *user_setting.OAuth2CommonHandlers {
	return &user_setting.OAuth2CommonHandlers{
		OwnerID:            org.Organization.ID,
		BasePathList:       fmt.Sprintf("%s/org/%s/settings/applications", setting.AppSubURL, org.Organization.Name),
		BasePathEditPrefix: fmt.Sprintf("%s/org/%s/settings/applications/oauth2", setting.AppSubURL, org.Organization.Name),
		TplAppEdit:         tplSettingsOAuthApplicationEdit,
	}
}

// Applications render org applications page (for org, at the moment, there are only OAuth2 applications)
func Applications(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.applications")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsApplications"] = true

	apps, err := db.Find[auth.OAuth2Application](ctx, auth.FindOAuth2ApplicationsOptions{
		OwnerID: ctx.Org.Organization.ID,
	})
	if err != nil {
		ctx.ServerError("GetOAuth2ApplicationsByUserID", err)
		return
	}
	ctx.Data["Applications"] = apps

	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	ctx.HTML(http.StatusOK, tplSettingsApplications)
}

// OAuthApplicationsPost response for adding an oauth2 application
func OAuthApplicationsPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.applications")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsApplications"] = true

	oa := newOAuth2CommonHandlers(ctx.Org)
	oa.AddApp(ctx)
}

// OAuth2ApplicationShow displays the given application
func OAuth2ApplicationShow(ctx *context.Context) {
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsApplications"] = true

	oa := newOAuth2CommonHandlers(ctx.Org)
	oa.EditShow(ctx)
}

// OAuth2ApplicationEdit response for editing oauth2 application
func OAuth2ApplicationEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.applications")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsApplications"] = true

	oa := newOAuth2CommonHandlers(ctx.Org)
	oa.EditSave(ctx)
}

// OAuthApplicationsRegenerateSecret handles the post request for regenerating the secret
func OAuthApplicationsRegenerateSecret(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsApplications"] = true

	oa := newOAuth2CommonHandlers(ctx.Org)
	oa.RegenerateSecret(ctx)
}

// DeleteOAuth2Application deletes the given oauth2 application
func DeleteOAuth2Application(ctx *context.Context) {
	oa := newOAuth2CommonHandlers(ctx.Org)
	oa.DeleteApp(ctx)
}

// TODO: revokes the grant with the given id
