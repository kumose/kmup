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

package admin

import (
	"net/http"

	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
)

// ListUserBadges lists all badges belonging to a user
func ListUserBadges(ctx *context.APIContext) {
	// swagger:operation GET /admin/users/{username}/badges admin adminListUserBadges
	// ---
	// summary: List a user's badges
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of the user whose badges are to be listed
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/BadgeList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	badges, maxResults, err := user_model.GetUserBadges(ctx, ctx.ContextUser)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.SetTotalCountHeader(maxResults)
	ctx.JSON(http.StatusOK, &badges)
}

// AddUserBadges add badges to a user
func AddUserBadges(ctx *context.APIContext) {
	// swagger:operation POST /admin/users/{username}/badges admin adminAddUserBadges
	// ---
	// summary: Add a badge to a user
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of the user to whom a badge is to be added
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/UserBadgeOption"
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "403":
	//     "$ref": "#/responses/forbidden"

	form := web.GetForm(ctx).(*api.UserBadgeOption)
	badges := prepareBadgesForReplaceOrAdd(*form)

	if err := user_model.AddUserBadges(ctx, ctx.ContextUser, badges); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteUserBadges delete a badge from a user
func DeleteUserBadges(ctx *context.APIContext) {
	// swagger:operation DELETE /admin/users/{username}/badges admin adminDeleteUserBadges
	// ---
	// summary: Remove a badge from a user
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of the user whose badge is to be deleted
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/UserBadgeOption"
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	//   "422":
	//     "$ref": "#/responses/validationError"

	form := web.GetForm(ctx).(*api.UserBadgeOption)
	badges := prepareBadgesForReplaceOrAdd(*form)

	if err := user_model.RemoveUserBadges(ctx, ctx.ContextUser, badges); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func prepareBadgesForReplaceOrAdd(form api.UserBadgeOption) []*user_model.Badge {
	badges := make([]*user_model.Badge, len(form.BadgeSlugs))
	for i, badge := range form.BadgeSlugs {
		badges[i] = &user_model.Badge{
			Slug: badge,
		}
	}
	return badges
}
