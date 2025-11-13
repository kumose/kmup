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

package shared

import (
	"errors"
	"net/http"

	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
	user_service "github.com/kumose/kmup/services/user"
)

func ListBlocks(ctx *context.APIContext, blocker *user_model.User) {
	blocks, total, err := user_model.FindBlockings(ctx, &user_model.FindBlockingOptions{
		ListOptions: utils.GetListOptions(ctx),
		BlockerID:   blocker.ID,
	})
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	if err := user_model.BlockingList(blocks).LoadAttributes(ctx); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	users := make([]*api.User, 0, len(blocks))
	for _, b := range blocks {
		users = append(users, convert.ToUser(ctx, b.Blockee, blocker))
	}

	ctx.SetTotalCountHeader(total)
	ctx.JSON(http.StatusOK, &users)
}

func CheckUserBlock(ctx *context.APIContext, blocker *user_model.User) {
	blockee, err := user_model.GetUserByName(ctx, ctx.PathParam("username"))
	if err != nil {
		ctx.APIErrorNotFound("GetUserByName", err)
		return
	}

	status := http.StatusNotFound
	blocking, err := user_model.GetBlocking(ctx, blocker.ID, blockee.ID)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	if blocking != nil {
		status = http.StatusNoContent
	}

	ctx.Status(status)
}

func BlockUser(ctx *context.APIContext, blocker *user_model.User) {
	blockee, err := user_model.GetUserByName(ctx, ctx.PathParam("username"))
	if err != nil {
		ctx.APIErrorNotFound("GetUserByName", err)
		return
	}

	if err := user_service.BlockUser(ctx, ctx.Doer, blocker, blockee, ctx.FormString("note")); err != nil {
		if errors.Is(err, user_model.ErrCanNotBlock) || errors.Is(err, user_model.ErrBlockOrganization) {
			ctx.APIError(http.StatusBadRequest, err)
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func UnblockUser(ctx *context.APIContext, doer, blocker *user_model.User) {
	blockee, err := user_model.GetUserByName(ctx, ctx.PathParam("username"))
	if err != nil {
		ctx.APIErrorNotFound("GetUserByName", err)
		return
	}

	if err := user_service.UnblockUser(ctx, doer, blocker, blockee); err != nil {
		if errors.Is(err, user_model.ErrCanNotUnblock) || errors.Is(err, user_model.ErrBlockOrganization) {
			ctx.APIError(http.StatusBadRequest, err)
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
