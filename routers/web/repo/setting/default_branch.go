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

	git_model "github.com/kumose/kmup/models/git"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/routers/web/repo"
	"github.com/kumose/kmup/services/context"
	repo_service "github.com/kumose/kmup/services/repository"
)

// SetDefaultBranchPost set default branch
func SetDefaultBranchPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.branches.update_default_branch")
	ctx.Data["PageIsSettingsBranches"] = true

	repo.PrepareBranchList(ctx)
	if ctx.Written() {
		return
	}

	repo := ctx.Repo.Repository

	switch ctx.FormString("action") {
	case "default_branch":
		if ctx.HasError() {
			ctx.HTML(http.StatusOK, tplBranches)
			return
		}

		branch := ctx.FormString("branch")
		if err := repo_service.SetRepoDefaultBranch(ctx, ctx.Repo.Repository, branch); err != nil {
			switch {
			case git_model.IsErrBranchNotExist(err):
				ctx.Status(http.StatusNotFound)
			default:
				ctx.ServerError("SetDefaultBranch", err)
			}
			return
		}

		log.Trace("Repository basic settings updated: %s/%s", ctx.Repo.Owner.Name, repo.Name)

		ctx.Flash.Success(ctx.Tr("repo.settings.update_settings_success"))
		ctx.Redirect(setting.AppSubURL + ctx.Req.URL.EscapedPath())
	default:
		ctx.NotFound(nil)
	}
}
