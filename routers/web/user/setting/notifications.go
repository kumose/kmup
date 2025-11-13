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

	"github.com/kumose/kmup/models/unit"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/user"
)

const tplSettingsNotifications templates.TplName = "user/settings/notifications"

// Notifications render user's notifications settings
func Notifications(ctx *context.Context) {
	if !setting.Service.EnableNotifyMail {
		ctx.NotFound(nil)
		return
	}

	ctx.Data["Title"] = ctx.Tr("notifications")
	ctx.Data["PageIsSettingsNotifications"] = true
	ctx.Data["EmailNotificationsPreference"] = ctx.Doer.EmailNotificationsPreference

	actionsEmailPref, err := user_model.GetUserSetting(ctx, ctx.Doer.ID, user_model.SettingsKeyEmailNotificationKmupActions, user_model.SettingEmailNotificationKmupActionsFailureOnly)
	if err != nil {
		ctx.ServerError("GetUserSetting", err)
		return
	}
	ctx.Data["ActionsEmailNotificationsPreference"] = actionsEmailPref

	ctx.HTML(http.StatusOK, tplSettingsNotifications)
}

// NotificationsEmailPost set user's email notification preference
func NotificationsEmailPost(ctx *context.Context) {
	if !setting.Service.EnableNotifyMail {
		ctx.NotFound(nil)
		return
	}

	preference := ctx.FormString("preference")
	if !(preference == user_model.EmailNotificationsEnabled ||
		preference == user_model.EmailNotificationsOnMention ||
		preference == user_model.EmailNotificationsDisabled ||
		preference == user_model.EmailNotificationsAndYourOwn) {
		ctx.Flash.Error(ctx.Tr("invalid_data", preference))
		ctx.Redirect(setting.AppSubURL + "/user/settings/notifications")
		return
	}
	opts := &user.UpdateOptions{
		EmailNotificationsPreference: optional.Some(preference),
	}
	if err := user.UpdateUser(ctx, ctx.Doer, opts); err != nil {
		ctx.ServerError("UpdateUser", err)
		return
	}
	ctx.Flash.Success(ctx.Tr("settings.email_preference_set_success"))
	ctx.Redirect(setting.AppSubURL + "/user/settings/notifications")
}

// NotificationsActionsEmailPost set user's email notification preference on Kmup Actions
func NotificationsActionsEmailPost(ctx *context.Context) {
	if !setting.Actions.Enabled || unit.TypeActions.UnitGlobalDisabled() {
		ctx.NotFound(nil)
		return
	}

	preference := ctx.FormString("preference")
	if !(preference == user_model.SettingEmailNotificationKmupActionsAll ||
		preference == user_model.SettingEmailNotificationKmupActionsDisabled ||
		preference == user_model.SettingEmailNotificationKmupActionsFailureOnly) {
		ctx.Flash.Error(ctx.Tr("invalid_data", preference))
		ctx.Redirect(setting.AppSubURL + "/user/settings/notifications")
		return
	}
	if err := user_model.SetUserSetting(ctx, ctx.Doer.ID, user_model.SettingsKeyEmailNotificationKmupActions, preference); err != nil {
		ctx.ServerError("SetUserSetting", err)
		return
	}
	ctx.Flash.Success(ctx.Tr("settings.email_preference_set_success"))
	ctx.Redirect(setting.AppSubURL + "/user/settings/notifications")
}
