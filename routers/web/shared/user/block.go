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
	"errors"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/forms"
	user_service "github.com/kumose/kmup/services/user"
)

func BlockedUsers(ctx *context.Context, blocker *user_model.User) {
	blocks, _, err := user_model.FindBlockings(ctx, &user_model.FindBlockingOptions{
		BlockerID: blocker.ID,
	})
	if err != nil {
		ctx.ServerError("FindBlockings", err)
		return
	}
	if err := user_model.BlockingList(blocks).LoadAttributes(ctx); err != nil {
		ctx.ServerError("LoadAttributes", err)
		return
	}
	ctx.Data["UserBlocks"] = blocks
}

func BlockedUsersPost(ctx *context.Context, blocker *user_model.User) {
	form := web.GetForm(ctx).(*forms.BlockUserForm)
	if ctx.HasError() {
		ctx.ServerError("FormValidation", nil)
		return
	}

	blockee, err := user_model.GetUserByName(ctx, form.Blockee)
	if err != nil {
		ctx.ServerError("GetUserByName", nil)
		return
	}

	switch form.Action {
	case "block":
		if err := user_service.BlockUser(ctx, ctx.Doer, blocker, blockee, form.Note); err != nil {
			if errors.Is(err, user_model.ErrCanNotBlock) || errors.Is(err, user_model.ErrBlockOrganization) {
				ctx.Flash.Error(ctx.Tr("user.block.block.failure", err.Error()))
			} else {
				ctx.ServerError("BlockUser", err)
				return
			}
		}
	case "unblock":
		if err := user_service.UnblockUser(ctx, ctx.Doer, blocker, blockee); err != nil {
			if errors.Is(err, user_model.ErrCanNotUnblock) || errors.Is(err, user_model.ErrBlockOrganization) {
				ctx.Flash.Error(ctx.Tr("user.block.unblock.failure", err.Error()))
			} else {
				ctx.ServerError("UnblockUser", err)
				return
			}
		}
	case "note":
		block, err := user_model.GetBlocking(ctx, blocker.ID, blockee.ID)
		if err != nil {
			ctx.ServerError("GetBlocking", err)
			return
		}
		if block != nil {
			if err := user_model.UpdateBlockingNote(ctx, block.ID, form.Note); err != nil {
				ctx.ServerError("UpdateBlockingNote", err)
				return
			}
		}
	}
}
