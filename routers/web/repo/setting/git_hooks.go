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

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/services/context"
)

// GitHooks hooks of a repository
func GitHooks(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.githooks")
	ctx.Data["PageIsSettingsGitHooks"] = true

	hooks, err := ctx.Repo.GitRepo.Hooks()
	if err != nil {
		ctx.ServerError("Hooks", err)
		return
	}
	ctx.Data["Hooks"] = hooks

	ctx.HTML(http.StatusOK, tplGithooks)
}

// GitHooksEdit render for editing a hook of repository page
func GitHooksEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.githooks")
	ctx.Data["PageIsSettingsGitHooks"] = true

	name := ctx.PathParam("name")
	hook, err := ctx.Repo.GitRepo.GetHook(name)
	if err != nil {
		if err == git.ErrNotValidHook {
			ctx.NotFound(err)
		} else {
			ctx.ServerError("GetHook", err)
		}
		return
	}
	ctx.Data["Hook"] = hook
	ctx.HTML(http.StatusOK, tplGithookEdit)
}

// GitHooksEditPost response for editing a git hook of a repository
func GitHooksEditPost(ctx *context.Context) {
	name := ctx.PathParam("name")
	hook, err := ctx.Repo.GitRepo.GetHook(name)
	if err != nil {
		if err == git.ErrNotValidHook {
			ctx.NotFound(err)
		} else {
			ctx.ServerError("GetHook", err)
		}
		return
	}
	hook.Content = ctx.FormString("content")
	if err = hook.Update(); err != nil {
		ctx.ServerError("hook.Update", err)
		return
	}
	ctx.Redirect(ctx.Repo.RepoLink + "/settings/hooks/git")
}
