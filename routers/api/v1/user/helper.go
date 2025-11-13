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
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/services/context"
)

// GetUserByPathParam get user by the path param name
// it will redirect to the user's new name if the user's name has been changed
func GetUserByPathParam(ctx *context.APIContext, name string) *user_model.User {
	username := ctx.PathParam(name)
	user, err := user_model.GetUserByName(ctx, username)
	if err != nil {
		if user_model.IsErrUserNotExist(err) {
			if redirectUserID, err2 := user_model.LookupUserRedirect(ctx, username); err2 == nil {
				context.RedirectToUser(ctx.Base, username, redirectUserID)
			} else {
				ctx.APIErrorNotFound("GetUserByName", err)
			}
		} else {
			ctx.APIErrorInternal(err)
		}
		return nil
	}
	return user
}

// GetContextUserByPathParam returns user whose name is presented in URL (path param "username").
func GetContextUserByPathParam(ctx *context.APIContext) *user_model.User {
	return GetUserByPathParam(ctx, "username")
}
