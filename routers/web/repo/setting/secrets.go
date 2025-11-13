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
	"errors"
	"net/http"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	shared "github.com/kumose/kmup/routers/web/shared/secrets"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	"github.com/kumose/kmup/services/context"
)

const (
	// TODO: Separate secrets from runners when layout is ready
	tplRepoSecrets templates.TplName = "repo/settings/actions"
	tplOrgSecrets  templates.TplName = "org/settings/actions"
	tplUserSecrets templates.TplName = "user/settings/actions"
)

type secretsCtx struct {
	OwnerID         int64
	RepoID          int64
	IsRepo          bool
	IsOrg           bool
	IsUser          bool
	SecretsTemplate templates.TplName
	RedirectLink    string
}

func getSecretsCtx(ctx *context.Context) (*secretsCtx, error) {
	if ctx.Data["PageIsRepoSettings"] == true {
		return &secretsCtx{
			OwnerID:         0,
			RepoID:          ctx.Repo.Repository.ID,
			IsRepo:          true,
			SecretsTemplate: tplRepoSecrets,
			RedirectLink:    ctx.Repo.RepoLink + "/settings/actions/secrets",
		}, nil
	}

	if ctx.Data["PageIsOrgSettings"] == true {
		if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
			ctx.ServerError("RenderUserOrgHeader", err)
			return nil, nil
		}
		return &secretsCtx{
			OwnerID:         ctx.ContextUser.ID,
			RepoID:          0,
			IsOrg:           true,
			SecretsTemplate: tplOrgSecrets,
			RedirectLink:    ctx.Org.OrgLink + "/settings/actions/secrets",
		}, nil
	}

	if ctx.Data["PageIsUserSettings"] == true {
		return &secretsCtx{
			OwnerID:         ctx.Doer.ID,
			RepoID:          0,
			IsUser:          true,
			SecretsTemplate: tplUserSecrets,
			RedirectLink:    setting.AppSubURL + "/user/settings/actions/secrets",
		}, nil
	}

	return nil, errors.New("unable to set Secrets context")
}

func Secrets(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("actions.actions")
	ctx.Data["PageType"] = "secrets"
	ctx.Data["PageIsSharedSettingsSecrets"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	sCtx, err := getSecretsCtx(ctx)
	if err != nil {
		ctx.ServerError("getSecretsCtx", err)
		return
	}

	if sCtx.IsRepo {
		ctx.Data["DisableSSH"] = setting.SSH.Disabled
	}

	shared.SetSecretsContext(ctx, sCtx.OwnerID, sCtx.RepoID)
	if ctx.Written() {
		return
	}
	ctx.HTML(http.StatusOK, sCtx.SecretsTemplate)
}

func SecretsPost(ctx *context.Context) {
	sCtx, err := getSecretsCtx(ctx)
	if err != nil {
		ctx.ServerError("getSecretsCtx", err)
		return
	}

	if ctx.HasError() {
		ctx.JSONError(ctx.GetErrMsg())
		return
	}

	shared.PerformSecretsPost(
		ctx,
		sCtx.OwnerID,
		sCtx.RepoID,
		sCtx.RedirectLink,
	)
}

func SecretsDelete(ctx *context.Context) {
	sCtx, err := getSecretsCtx(ctx)
	if err != nil {
		ctx.ServerError("getSecretsCtx", err)
		return
	}
	shared.PerformSecretsDelete(
		ctx,
		sCtx.OwnerID,
		sCtx.RepoID,
		sCtx.RedirectLink,
	)
}
