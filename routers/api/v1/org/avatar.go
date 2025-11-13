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

package org

import (
	"encoding/base64"
	"net/http"

	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	user_service "github.com/kumose/kmup/services/user"
)

// UpdateAvatarupdates the Avatar of an Organisation
func UpdateAvatar(ctx *context.APIContext) {
	// swagger:operation POST /orgs/{org}/avatar organization orgUpdateAvatar
	// ---
	// summary: Update Avatar
	// produces:
	// - application/json
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/UpdateUserAvatarOption"
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	form := web.GetForm(ctx).(*api.UpdateUserAvatarOption)

	content, err := base64.StdEncoding.DecodeString(form.Image)
	if err != nil {
		ctx.APIError(http.StatusBadRequest, err)
		return
	}

	err = user_service.UploadAvatar(ctx, ctx.Org.Organization.AsUser(), content)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteAvatar deletes the Avatar of an Organisation
func DeleteAvatar(ctx *context.APIContext) {
	// swagger:operation DELETE /orgs/{org}/avatar organization orgDeleteAvatar
	// ---
	// summary: Delete Avatar
	// produces:
	// - application/json
	// parameters:
	// - name: org
	//   in: path
	//   description: name of the organization
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	err := user_service.DeleteAvatar(ctx, ctx.Org.Organization.AsUser())
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
