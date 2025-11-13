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

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	shared "github.com/kumose/kmup/routers/web/shared/packages"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	"github.com/kumose/kmup/services/context"
)

const (
	tplSettingsPackages            templates.TplName = "org/settings/packages"
	tplSettingsPackagesRuleEdit    templates.TplName = "org/settings/packages_cleanup_rules_edit"
	tplSettingsPackagesRulePreview templates.TplName = "org/settings/packages_cleanup_rules_preview"
)

func Packages(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	shared.SetPackagesContext(ctx, ctx.ContextUser)

	ctx.HTML(http.StatusOK, tplSettingsPackages)
}

func PackagesRuleAdd(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	shared.SetRuleAddContext(ctx)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRuleEdit)
}

func PackagesRuleEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	shared.SetRuleEditContext(ctx, ctx.ContextUser)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRuleEdit)
}

func PackagesRuleAddPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	shared.PerformRuleAddPost(
		ctx,
		ctx.ContextUser,
		fmt.Sprintf("%s/org/%s/settings/packages", setting.AppSubURL, ctx.ContextUser.Name),
		tplSettingsPackagesRuleEdit,
	)
}

func PackagesRuleEditPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	shared.PerformRuleEditPost(
		ctx,
		ctx.ContextUser,
		fmt.Sprintf("%s/org/%s/settings/packages", setting.AppSubURL, ctx.ContextUser.Name),
		tplSettingsPackagesRuleEdit,
	)
}

func PackagesRulePreview(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	shared.SetRulePreviewContext(ctx, ctx.ContextUser)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRulePreview)
}

func InitializeCargoIndex(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	shared.InitializeCargoIndex(ctx, ctx.ContextUser)

	ctx.Redirect(fmt.Sprintf("%s/org/%s/settings/packages", setting.AppSubURL, ctx.ContextUser.Name))
}

func RebuildCargoIndex(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsPackages"] = true

	shared.RebuildCargoIndex(ctx, ctx.ContextUser)

	ctx.Redirect(fmt.Sprintf("%s/org/%s/settings/packages", setting.AppSubURL, ctx.ContextUser.Name))
}
