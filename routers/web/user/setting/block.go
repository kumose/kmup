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

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	"github.com/kumose/kmup/services/context"
)

const (
	tplSettingsBlockedUsers templates.TplName = "user/settings/blocked_users"
)

func BlockedUsers(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("user.block.list")
	ctx.Data["PageIsSettingsBlockedUsers"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared_user.BlockedUsers(ctx, ctx.Doer)
	if ctx.Written() {
		return
	}

	ctx.HTML(http.StatusOK, tplSettingsBlockedUsers)
}

func BlockedUsersPost(ctx *context.Context) {
	shared_user.BlockedUsersPost(ctx, ctx.Doer)
	if ctx.Written() {
		return
	}

	ctx.Redirect(setting.AppSubURL + "/user/settings/blocked_users")
}
