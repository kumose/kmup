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

package context

import (
	"fmt"
	"net/http"
	"strings"

	user_model "github.com/kumose/kmup/models/user"
)

// UserAssignmentWeb returns a middleware to handle context-user assignment for web routes
func UserAssignmentWeb() func(ctx *Context) {
	return func(ctx *Context) {
		errorFn := func(status int, obj any) {
			err, ok := obj.(error)
			if !ok {
				err = fmt.Errorf("%s", obj)
			}
			if status == http.StatusNotFound {
				ctx.NotFound(err)
			} else {
				ctx.ServerError("UserAssignmentWeb", err)
			}
		}
		ctx.ContextUser = userAssignment(ctx.Base, ctx.Doer, errorFn)
		ctx.Data["ContextUser"] = ctx.ContextUser
	}
}

// UserIDAssignmentAPI returns a middleware to handle context-user assignment for api routes
func UserIDAssignmentAPI() func(ctx *APIContext) {
	return func(ctx *APIContext) {
		userID := ctx.PathParamInt64("user-id")

		if ctx.IsSigned && ctx.Doer.ID == userID {
			ctx.ContextUser = ctx.Doer
		} else {
			var err error
			ctx.ContextUser, err = user_model.GetUserByID(ctx, userID)
			if err != nil {
				if user_model.IsErrUserNotExist(err) {
					ctx.APIError(http.StatusNotFound, err)
				} else {
					ctx.APIErrorInternal(err)
				}
			}
		}
	}
}

// UserAssignmentAPI returns a middleware to handle context-user assignment for api routes
func UserAssignmentAPI() func(ctx *APIContext) {
	return func(ctx *APIContext) {
		ctx.ContextUser = userAssignment(ctx.Base, ctx.Doer, ctx.APIError)
	}
}

func userAssignment(ctx *Base, doer *user_model.User, errCb func(int, any)) (contextUser *user_model.User) {
	username := ctx.PathParam("username")

	if doer != nil && strings.EqualFold(doer.LowerName, username) {
		contextUser = doer
	} else {
		var err error
		contextUser, err = user_model.GetUserByName(ctx, username)
		if err != nil {
			if user_model.IsErrUserNotExist(err) {
				if redirectUserID, err := user_model.LookupUserRedirect(ctx, username); err == nil {
					RedirectToUser(ctx, username, redirectUserID)
				} else if user_model.IsErrUserRedirectNotExist(err) {
					errCb(http.StatusNotFound, err)
				} else {
					errCb(http.StatusInternalServerError, fmt.Errorf("LookupUserRedirect: %w", err))
				}
			} else {
				errCb(http.StatusInternalServerError, fmt.Errorf("GetUserByName: %w", err))
			}
		}
	}
	return contextUser
}
