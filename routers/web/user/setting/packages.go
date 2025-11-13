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
	"strings"

	user_model "github.com/kumose/kmup/models/user"
	chef_module "github.com/kumose/kmup/modules/packages/chef"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/util"
	shared "github.com/kumose/kmup/routers/web/shared/packages"
	"github.com/kumose/kmup/services/context"
)

const (
	tplSettingsPackages            templates.TplName = "user/settings/packages"
	tplSettingsPackagesRuleEdit    templates.TplName = "user/settings/packages_cleanup_rules_edit"
	tplSettingsPackagesRulePreview templates.TplName = "user/settings/packages_cleanup_rules_preview"
)

func Packages(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetPackagesContext(ctx, ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsPackages)
}

func PackagesRuleAdd(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetRuleAddContext(ctx)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRuleEdit)
}

func PackagesRuleEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetRuleEditContext(ctx, ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRuleEdit)
}

func PackagesRuleAddPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.PerformRuleAddPost(
		ctx,
		ctx.Doer,
		setting.AppSubURL+"/user/settings/packages",
		tplSettingsPackagesRuleEdit,
	)
}

func PackagesRuleEditPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.PerformRuleEditPost(
		ctx,
		ctx.Doer,
		setting.AppSubURL+"/user/settings/packages",
		tplSettingsPackagesRuleEdit,
	)
}

func PackagesRulePreview(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetRulePreviewContext(ctx, ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRulePreview)
}

func InitializeCargoIndex(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true

	shared.InitializeCargoIndex(ctx, ctx.Doer)

	ctx.Redirect(setting.AppSubURL + "/user/settings/packages")
}

func RebuildCargoIndex(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true

	shared.RebuildCargoIndex(ctx, ctx.Doer)

	ctx.Redirect(setting.AppSubURL + "/user/settings/packages")
}

func RegenerateChefKeyPair(ctx *context.Context) {
	priv, pub, err := util.GenerateKeyPair(chef_module.KeyBits)
	if err != nil {
		ctx.ServerError("GenerateKeyPair", err)
		return
	}

	if err := user_model.SetUserSetting(ctx, ctx.Doer.ID, chef_module.SettingPublicPem, pub); err != nil {
		ctx.ServerError("SetUserSetting", err)
		return
	}

	ctx.ServeContent(strings.NewReader(priv), &context.ServeHeaderOptions{
		ContentType: "application/x-pem-file",
		Filename:    ctx.Doer.Name + ".priv",
	})
}
