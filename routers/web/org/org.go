// Copyright 2014 The Gogs Authors. All rights reserved.
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
	"errors"
	"net/http"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/organization"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/forms"
)

const (
	// tplCreateOrg template path for create organization
	tplCreateOrg templates.TplName = "org/create"
)

// Create render the page for create organization
func Create(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("new_org")
	if !ctx.Doer.CanCreateOrganization() {
		ctx.ServerError("Not allowed", errors.New(ctx.Locale.TrString("org.form.create_org_not_allowed")))
		return
	}

	ctx.Data["visibility"] = setting.Service.DefaultOrgVisibilityMode
	ctx.Data["repo_admin_change_team_access"] = true

	ctx.HTML(http.StatusOK, tplCreateOrg)
}

// CreatePost response for create organization
func CreatePost(ctx *context.Context) {
	form := *web.GetForm(ctx).(*forms.CreateOrgForm)
	ctx.Data["Title"] = ctx.Tr("new_org")

	if !ctx.Doer.CanCreateOrganization() {
		ctx.ServerError("Not allowed", errors.New(ctx.Locale.TrString("org.form.create_org_not_allowed")))
		return
	}

	if ctx.HasError() {
		ctx.HTML(http.StatusOK, tplCreateOrg)
		return
	}

	org := &organization.Organization{
		Name:                      form.OrgName,
		IsActive:                  true,
		Type:                      user_model.UserTypeOrganization,
		Visibility:                form.Visibility,
		RepoAdminChangeTeamAccess: form.RepoAdminChangeTeamAccess,
	}

	if err := organization.CreateOrganization(ctx, org, ctx.Doer); err != nil {
		ctx.Data["Err_OrgName"] = true
		switch {
		case user_model.IsErrUserAlreadyExist(err):
			ctx.RenderWithErr(ctx.Tr("form.org_name_been_taken"), tplCreateOrg, &form)
		case db.IsErrNameReserved(err):
			ctx.RenderWithErr(ctx.Tr("org.form.name_reserved", err.(db.ErrNameReserved).Name), tplCreateOrg, &form)
		case db.IsErrNamePatternNotAllowed(err):
			ctx.RenderWithErr(ctx.Tr("org.form.name_pattern_not_allowed", err.(db.ErrNamePatternNotAllowed).Pattern), tplCreateOrg, &form)
		case organization.IsErrUserNotAllowedCreateOrg(err):
			ctx.RenderWithErr(ctx.Tr("org.form.create_org_not_allowed"), tplCreateOrg, &form)
		default:
			ctx.ServerError("CreateOrganization", err)
		}
		return
	}
	log.Trace("Organization created: %s", org.Name)

	ctx.Redirect(org.AsUser().DashboardLink())
}
