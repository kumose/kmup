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
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// ListStargazers list a repository's stargazers
func ListStargazers(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/stargazers repository repoListStargazers
	// ---
	// summary: List a repo's stargazers
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
	//     "$ref": "#/responses/UserList"
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "403":
	//     "$ref": "#/responses/forbidden"

	stargazers, err := repo_model.GetStargazers(ctx, ctx.Repo.Repository, utils.GetListOptions(ctx))
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	users := make([]*api.User, len(stargazers))
	for i, stargazer := range stargazers {
		users[i] = convert.ToUser(ctx, stargazer, ctx.Doer)
	}

	ctx.SetTotalCountHeader(int64(ctx.Repo.Repository.NumStars))
	ctx.JSON(http.StatusOK, users)
}
