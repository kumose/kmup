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
	"time"

	"github.com/kumose/kmup/models/avatars"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/httpcache"
	"github.com/kumose/kmup/services/context"
)

func cacheableRedirect(ctx *context.Context, location string) {
	// here we should not use `setting.StaticCacheTime`, it is pretty long (default: 6 hours)
	// we must make sure the redirection cache time is short enough, otherwise a user won't see the updated avatar in 6 hours
	// it's OK to make the cache time short, it is only a redirection, and doesn't cost much to make a new request
	httpcache.SetCacheControlInHeader(ctx.Resp.Header(), &httpcache.CacheControlOptions{MaxAge: 5 * time.Minute})
	ctx.Redirect(location)
}

// AvatarByUsernameSize redirect browser to user avatar of requested size
func AvatarByUsernameSize(ctx *context.Context) {
	username := ctx.PathParam("username")
	user := user_model.GetSystemUserByName(username)
	if user == nil {
		var err error
		if user, err = user_model.GetUserByName(ctx, username); err != nil {
			ctx.NotFoundOrServerError("GetUserByName", user_model.IsErrUserNotExist, err)
			return
		}
	}
	cacheableRedirect(ctx, user.AvatarLinkWithSize(ctx, ctx.PathParamInt("size")))
}

// AvatarByEmailHash redirects the browser to the email avatar link
func AvatarByEmailHash(ctx *context.Context) {
	hash := ctx.PathParam("hash")
	email, err := avatars.GetEmailForHash(ctx, hash)
	if err != nil {
		ctx.ServerError("invalid avatar hash: "+hash, err)
		return
	}
	size := ctx.FormInt("size")
	cacheableRedirect(ctx, avatars.GenerateEmailAvatarFinalLink(ctx, email, size))
}
