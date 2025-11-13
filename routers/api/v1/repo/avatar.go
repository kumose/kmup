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
	"encoding/base64"
	"net/http"

	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	repo_service "github.com/kumose/kmup/services/repository"
)

// UpdateVatar updates the Avatar of an Repo
func UpdateAvatar(ctx *context.APIContext) {
	// swagger:operation POST /repos/{owner}/{repo}/avatar repository repoUpdateAvatar
	// ---
	// summary: Update avatar
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
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/UpdateRepoAvatarOption"
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	form := web.GetForm(ctx).(*api.UpdateRepoAvatarOption)

	content, err := base64.StdEncoding.DecodeString(form.Image)
	if err != nil {
		ctx.APIError(http.StatusBadRequest, err)
		return
	}

	err = repo_service.UploadAvatar(ctx, ctx.Repo.Repository, content)
	if err != nil {
		ctx.APIErrorInternal(err)
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateAvatar deletes the Avatar of an Repo
func DeleteAvatar(ctx *context.APIContext) {
	// swagger:operation DELETE /repos/{owner}/{repo}/avatar repository repoDeleteAvatar
	// ---
	// summary: Delete avatar
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
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	err := repo_service.DeleteAvatar(ctx, ctx.Repo.Repository)
	if err != nil {
		ctx.APIErrorInternal(err)
	}

	ctx.Status(http.StatusNoContent)
}
