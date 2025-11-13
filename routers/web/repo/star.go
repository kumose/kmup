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

package repo

import (
	"net/http"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
)

const tplStarUnstar templates.TplName = "repo/star_unstar"

func ActionStar(ctx *context.Context) {
	err := repo_model.StarRepo(ctx, ctx.Doer, ctx.Repo.Repository, ctx.PathParam("action") == "star")
	if err != nil {
		handleActionError(ctx, err)
		return
	}

	ctx.Data["IsStaringRepo"] = repo_model.IsStaring(ctx, ctx.Doer.ID, ctx.Repo.Repository.ID)
	ctx.Data["Repository"], err = repo_model.GetRepositoryByName(ctx, ctx.Repo.Repository.OwnerID, ctx.Repo.Repository.Name)
	if err != nil {
		ctx.ServerError("GetRepositoryByName", err)
		return
	}
	ctx.RespHeader().Add("hx-trigger", "refreshUserCards") // see the `hx-trigger="refreshUserCards ..."` comments in tmpl
	ctx.HTML(http.StatusOK, tplStarUnstar)
}
