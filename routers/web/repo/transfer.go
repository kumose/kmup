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
	"github.com/kumose/kmup/services/context"
	repo_service "github.com/kumose/kmup/services/repository"
)

func acceptTransfer(ctx *context.Context) {
	err := repo_service.AcceptTransferOwnership(ctx, ctx.Repo.Repository, ctx.Doer)
	if err == nil {
		ctx.Flash.Success(ctx.Tr("repo.settings.transfer.success"))
		ctx.Redirect(ctx.Repo.Repository.Link())
		return
	}
	handleActionError(ctx, err)
}

func rejectTransfer(ctx *context.Context) {
	err := repo_service.RejectRepositoryTransfer(ctx, ctx.Repo.Repository, ctx.Doer)
	if err == nil {
		ctx.Flash.Success(ctx.Tr("repo.settings.transfer.rejected"))
		ctx.Redirect(ctx.Repo.Repository.Link())
		return
	}
	handleActionError(ctx, err)
}

func ActionTransfer(ctx *context.Context) {
	switch ctx.PathParam("action") {
	case "accept_transfer":
		acceptTransfer(ctx)
	case "reject_transfer":
		rejectTransfer(ctx)
	}
}
