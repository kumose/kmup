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

	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
	repo_service "github.com/kumose/kmup/services/repository"
)

const tplEditorFork templates.TplName = "repo/editor/fork"

func ForkToEdit(ctx *context.Context) {
	ctx.HTML(http.StatusOK, tplEditorFork)
}

func ForkToEditPost(ctx *context.Context) {
	ForkRepoTo(ctx, ctx.Doer, repo_service.ForkRepoOptions{
		BaseRepo:     ctx.Repo.Repository,
		Name:         getUniqueRepositoryName(ctx, ctx.Doer.ID, ctx.Repo.Repository.Name),
		Description:  ctx.Repo.Repository.Description,
		SingleBranch: ctx.Repo.Repository.DefaultBranch, // maybe we only need the default branch in the fork?
	})
	if ctx.Written() {
		return
	}
	ctx.JSONRedirect("") // reload the page, the new fork should be editable now
}
