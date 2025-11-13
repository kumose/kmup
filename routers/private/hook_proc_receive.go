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

package private

import (
	"errors"
	"net/http"

	issues_model "github.com/kumose/kmup/models/issues"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/private"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/agit"
	kmup_context "github.com/kumose/kmup/services/context"
)

// HookProcReceive proc-receive hook - only handles agit Proc-Receive requests at present
func HookProcReceive(ctx *kmup_context.PrivateContext) {
	opts := web.GetForm(ctx).(*private.HookOptions)
	if !git.DefaultFeatures().SupportProcReceive {
		ctx.Status(http.StatusNotFound)
		return
	}

	results, err := agit.ProcReceive(ctx, ctx.Repo.Repository, ctx.Repo.GitRepo, opts)
	if err != nil {
		if errors.Is(err, issues_model.ErrMustCollaborator) {
			ctx.JSON(http.StatusUnauthorized, private.Response{
				Err: err.Error(), UserMsg: "You must be a collaborator to create pull request.",
			})
		} else if errors.Is(err, user_model.ErrBlockedUser) {
			ctx.JSON(http.StatusUnauthorized, private.Response{
				Err: err.Error(), UserMsg: "Cannot create pull request because you are blocked by the repository owner.",
			})
		} else {
			log.Error("agit.ProcReceive failed: %v", err)
			ctx.JSON(http.StatusInternalServerError, private.Response{
				Err: err.Error(),
			})
		}

		return
	}

	ctx.JSON(http.StatusOK, private.HookProcReceiveResult{
		Results: results,
	})
}
