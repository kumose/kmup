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

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// SearchCandidates searches candidate users for dropdown list
func SearchCandidates(ctx *context.Context) {
	searchUserTypes := []user_model.UserType{user_model.UserTypeIndividual}
	if ctx.FormBool("orgs") {
		searchUserTypes = append(searchUserTypes, user_model.UserTypeOrganization)
	}
	users, _, err := user_model.SearchUsers(ctx, user_model.SearchUserOptions{
		Actor:       ctx.Doer,
		Keyword:     ctx.FormTrim("q"),
		Types:       searchUserTypes,
		IsActive:    optional.Some(true),
		ListOptions: db.ListOptions{PageSize: setting.UI.MembersPagingNum},
	})
	if err != nil {
		ctx.ServerError("Unable to search users", err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"data": convert.ToUsers(ctx, ctx.Doer, users)})
}
