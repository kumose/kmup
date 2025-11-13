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
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
	release_service "github.com/kumose/kmup/services/release"
)

// GetReleaseByTag get a single release of a repository by tag name
func GetReleaseByTag(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/releases/tags/{tag} repository repoGetReleaseByTag
	// ---
	// summary: Get a release by tag name
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
	// - name: tag
	//   in: path
	//   description: tag name of the release to get
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/Release"
	//   "404":
	//     "$ref": "#/responses/notFound"

	tag := ctx.PathParam("tag")

	release, err := repo_model.GetRelease(ctx, ctx.Repo.Repository.ID, tag)
	if err != nil {
		if repo_model.IsErrReleaseNotExist(err) {
			ctx.APIErrorNotFound()
			return
		}
		ctx.APIErrorInternal(err)
		return
	}

	if release.IsTag {
		ctx.APIErrorNotFound()
		return
	}

	if err = release.LoadAttributes(ctx); err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	ctx.JSON(http.StatusOK, convert.ToAPIRelease(ctx, ctx.Repo.Repository, release))
}

// DeleteReleaseByTag delete a release from a repository by tag name
func DeleteReleaseByTag(ctx *context.APIContext) {
	// swagger:operation DELETE /repos/{owner}/{repo}/releases/tags/{tag} repository repoDeleteReleaseByTag
	// ---
	// summary: Delete a release by tag name
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
	// - name: tag
	//   in: path
	//   description: tag name of the release to delete
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "422":
	//     "$ref": "#/responses/validationError"

	tag := ctx.PathParam("tag")

	release, err := repo_model.GetRelease(ctx, ctx.Repo.Repository.ID, tag)
	if err != nil {
		if repo_model.IsErrReleaseNotExist(err) {
			ctx.APIErrorNotFound()
			return
		}
		ctx.APIErrorInternal(err)
		return
	}

	if release.IsTag {
		ctx.APIErrorNotFound()
		return
	}

	if err = release_service.DeleteReleaseByID(ctx, ctx.Repo.Repository, release, ctx.Doer, false); err != nil {
		if release_service.IsErrProtectedTagName(err) {
			ctx.APIError(http.StatusUnprocessableEntity, "user not allowed to delete protected tag")
			return
		}
		ctx.APIErrorInternal(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
