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

	"github.com/kumose/kmup/modules/optional"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
	user_service "github.com/kumose/kmup/services/user"
)

// GetUserSettings returns user settings
func GetUserSettings(ctx *context.APIContext) {
	// swagger:operation GET /user/settings user getUserSettings
	// ---
	// summary: Get user settings
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/UserSettings"
	ctx.JSON(http.StatusOK, convert.User2UserSettings(ctx.Doer))
}

// UpdateUserSettings returns user settings
func UpdateUserSettings(ctx *context.APIContext) {
	// swagger:operation PATCH /user/settings user updateUserSettings
	// ---
	// summary: Update user settings
	// parameters:
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/UserSettingsOptions"
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/UserSettings"

	form := web.GetForm(ctx).(*api.UserSettingsOptions)

	opts := &user_service.UpdateOptions{
		FullName:            optional.FromPtr(form.FullName),
		Description:         optional.FromPtr(form.Description),
		Website:             optional.FromPtr(form.Website),
		Location:            optional.FromPtr(form.Location),
		Language:            optional.FromPtr(form.Language),
		Theme:               optional.FromPtr(form.Theme),
		DiffViewStyle:       optional.FromPtr(form.DiffViewStyle),
		KeepEmailPrivate:    optional.FromPtr(form.HideEmail),
		KeepActivityPrivate: optional.FromPtr(form.HideActivity),
	}
	if err := user_service.UpdateUser(ctx, ctx.Doer, opts); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.JSON(http.StatusOK, convert.User2UserSettings(ctx.Doer))
}
