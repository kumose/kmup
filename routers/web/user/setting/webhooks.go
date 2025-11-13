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
	"net/http"

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/models/webhook"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
)

const (
	tplSettingsHooks templates.TplName = "user/settings/hooks"
)

// Webhooks render webhook list page
func Webhooks(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsHooks"] = true
	ctx.Data["BaseLink"] = setting.AppSubURL + "/user/settings/hooks"
	ctx.Data["BaseLinkNew"] = setting.AppSubURL + "/user/settings/hooks"
	ctx.Data["Description"] = ctx.Tr("settings.hooks.desc")
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	ws, err := db.Find[webhook.Webhook](ctx, webhook.ListWebhookOptions{OwnerID: ctx.Doer.ID})
	if err != nil {
		ctx.ServerError("ListWebhooksByOpts", err)
		return
	}

	ctx.Data["Webhooks"] = ws
	ctx.HTML(http.StatusOK, tplSettingsHooks)
}

// DeleteWebhook response for delete webhook
func DeleteWebhook(ctx *context.Context) {
	if err := webhook.DeleteWebhookByOwnerID(ctx, ctx.Doer.ID, ctx.FormInt64("id")); err != nil {
		ctx.Flash.Error("DeleteWebhookByOwnerID: " + err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("repo.settings.webhook_deletion_success"))
	}

	ctx.JSONRedirect(setting.AppSubURL + "/user/settings/hooks")
}
