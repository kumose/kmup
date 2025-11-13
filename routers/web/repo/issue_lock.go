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
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/forms"
)

// LockIssue locks an issue. This would limit commenting abilities to
// users with write access to the repo.
func LockIssue(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.IssueLockForm)
	issue := GetActionIssue(ctx)
	if ctx.Written() {
		return
	}

	if issue.IsLocked {
		ctx.JSONError(ctx.Tr("repo.issues.lock_duplicate"))
		return
	}

	if err := issues_model.LockIssue(ctx, &issues_model.IssueLockOptions{
		Doer:   ctx.Doer,
		Issue:  issue,
		Reason: form.Reason,
	}); err != nil {
		ctx.ServerError("LockIssue", err)
		return
	}

	ctx.JSONRedirect(issue.Link())
}

// UnlockIssue unlocks a previously locked issue.
func UnlockIssue(ctx *context.Context) {
	issue := GetActionIssue(ctx)
	if ctx.Written() {
		return
	}

	if !issue.IsLocked {
		ctx.JSONError(ctx.Tr("repo.issues.unlock_error"))
		return
	}

	if err := issues_model.UnlockIssue(ctx, &issues_model.IssueLockOptions{
		Doer:  ctx.Doer,
		Issue: issue,
	}); err != nil {
		ctx.ServerError("UnlockIssue", err)
		return
	}

	ctx.JSONRedirect(issue.Link())
}
