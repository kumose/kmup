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

package user

import (
	"net/http"

	access_model "github.com/kumose/kmup/models/perm/access"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// listUserRepos - List the repositories owned by the given user.
func listUserRepos(ctx *context.APIContext, u *user_model.User, private bool) {
	opts := utils.GetListOptions(ctx)

	repos, count, err := repo_model.GetUserRepositories(ctx, repo_model.SearchRepoOptions{
		Actor:       u,
		Private:     private,
		ListOptions: opts,
		OrderBy:     "id ASC",
	})
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	if err := repos.LoadAttributes(ctx); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	apiRepos := make([]*api.Repository, 0, len(repos))
	for i := range repos {
		permission, err := access_model.GetUserRepoPermission(ctx, repos[i], ctx.Doer)
		if err != nil {
			ctx.APIErrorInternal(err)
			return
		}
		if ctx.IsSigned && ctx.Doer.IsAdmin || permission.HasAnyUnitAccess() {
			apiRepos = append(apiRepos, convert.ToRepo(ctx, repos[i], permission))
		}
	}

	ctx.SetLinkHeader(int(count), opts.PageSize)
	ctx.SetTotalCountHeader(count)
	ctx.JSON(http.StatusOK, &apiRepos)
}

// ListUserRepos - list the repos owned by the given user.
func ListUserRepos(ctx *context.APIContext) {
	// swagger:operation GET /users/{username}/repos user userListRepos
	// ---
	// summary: List the repos owned by the given user
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of the user whose owned repos are to be listed
	//   type: string
	//   required: true
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/RepositoryList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	private := ctx.IsSigned
	listUserRepos(ctx, ctx.ContextUser, private)
}

// ListMyRepos - list the repositories you own or have access to.
func ListMyRepos(ctx *context.APIContext) {
	// swagger:operation GET /user/repos user userCurrentListRepos
	// ---
	// summary: List the repos that the authenticated user owns
	// produces:
	// - application/json
	// parameters:
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/RepositoryList"

	opts := repo_model.SearchRepoOptions{
		ListOptions:        utils.GetListOptions(ctx),
		Actor:              ctx.Doer,
		OwnerID:            ctx.Doer.ID,
		Private:            ctx.IsSigned,
		IncludeDescription: true,
	}

	repos, count, err := repo_model.SearchRepository(ctx, opts)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	results := make([]*api.Repository, len(repos))
	for i, repo := range repos {
		if err = repo.LoadOwner(ctx); err != nil {
			ctx.APIErrorInternal(err)
			return
		}
		permission, err := access_model.GetUserRepoPermission(ctx, repo, ctx.Doer)
		if err != nil {
			ctx.APIErrorInternal(err)
		}
		results[i] = convert.ToRepo(ctx, repo, permission)
	}

	ctx.SetLinkHeader(int(count), opts.ListOptions.PageSize)
	ctx.SetTotalCountHeader(count)
	ctx.JSON(http.StatusOK, &results)
}

// ListOrgRepos - list the repositories of an organization.
func ListOrgRepos(ctx *context.APIContext) {
	// swagger:operation GET /orgs/{org}/repos organization orgListRepos
	// ---
	// summary: List an organization's repos
	// produces:
	// - application/json
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/RepositoryList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	listUserRepos(ctx, ctx.Org.Organization.AsUser(), ctx.IsSigned)
}
