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
	"errors"
	"net/http"

	"github.com/kumose/kmup/modules/git"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// ListGitHooks list all Git hooks of a repository
func ListGitHooks(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/hooks/git repository repoListGitHooks
	// ---
	// summary: List the Git hooks in a repository
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/GitHookList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	hooks, err := ctx.Repo.GitRepo.Hooks()
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	apiHooks := make([]*api.GitHook, len(hooks))
	for i := range hooks {
		apiHooks[i] = convert.ToGitHook(hooks[i])
	}
	ctx.JSON(http.StatusOK, &apiHooks)
}

// GetGitHook get a repo's Git hook by id
func GetGitHook(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/hooks/git/{id} repository repoGetGitHook
	// ---
	// summary: Get a Git hook
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: id
	//   in: path
	//   description: id of the hook to get
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/GitHook"
	//   "404":
	//     "$ref": "#/responses/notFound"

	hookID := ctx.PathParam("id")
	hook, err := ctx.Repo.GitRepo.GetHook(hookID)
	if err != nil {
		if errors.Is(err, git.ErrNotValidHook) {
			ctx.APIErrorNotFound()
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}
	ctx.JSON(http.StatusOK, convert.ToGitHook(hook))
}

// EditGitHook modify a Git hook of a repository
func EditGitHook(ctx *context.APIContext) {
	// swagger:operation PATCH /repos/{owner}/{repo}/hooks/git/{id} repository repoEditGitHook
	// ---
	// summary: Edit a Git hook in a repository
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: id
	//   in: path
	//   description: id of the hook to get
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/EditGitHookOption"
	// responses:
	//   "200":
	//     "$ref": "#/responses/GitHook"
	//   "404":
	//     "$ref": "#/responses/notFound"

	form := web.GetForm(ctx).(*api.EditGitHookOption)
	hookID := ctx.PathParam("id")
	hook, err := ctx.Repo.GitRepo.GetHook(hookID)
	if err != nil {
		if errors.Is(err, git.ErrNotValidHook) {
			ctx.APIErrorNotFound()
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

	hook.Content = form.Content
	if err = hook.Update(); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.JSON(http.StatusOK, convert.ToGitHook(hook))
}

// DeleteGitHook delete a Git hook of a repository
func DeleteGitHook(ctx *context.APIContext) {
	// swagger:operation DELETE /repos/{owner}/{repo}/hooks/git/{id} repository repoDeleteGitHook
	// ---
	// summary: Delete a Git hook in a repository
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: id
	//   in: path
	//   description: id of the hook to get
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"

	hookID := ctx.PathParam("id")
	hook, err := ctx.Repo.GitRepo.GetHook(hookID)
	if err != nil {
		if errors.Is(err, git.ErrNotValidHook) {
			ctx.APIErrorNotFound()
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

	hook.Content = ""
	if err = hook.Update(); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
