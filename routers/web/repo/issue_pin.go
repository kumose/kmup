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
	"net/http"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/context"
)

// IssuePinOrUnpin pin or unpin a Issue
func IssuePinOrUnpin(ctx *context.Context) {
	issue := GetActionIssue(ctx)
	if ctx.Written() {
		return
	}

	// If we don't do this, it will crash when trying to add the pin event to the comment history
	err := issue.LoadRepo(ctx)
	if err != nil {
		ctx.ServerError("LoadRepo", err)
		return
	}

	// PinOrUnpin pins or unpins a Issue
	_, err = issues_model.GetIssuePin(ctx, issue)
	if err != nil && !db.IsErrNotExist(err) {
		ctx.ServerError("GetIssuePin", err)
		return
	}

	if db.IsErrNotExist(err) {
		err = issues_model.PinIssue(ctx, issue, ctx.Doer)
	} else {
		err = issues_model.UnpinIssue(ctx, issue, ctx.Doer)
	}

	if err != nil {
		if issues_model.IsErrIssueMaxPinReached(err) {
			ctx.JSONError(ctx.Tr("repo.issues.max_pinned"))
		} else {
			ctx.ServerError("Pin/Unpin failed", err)
		}
		return
	}

	ctx.JSONRedirect(issue.Link())
}

// IssueUnpin unpins a Issue
func IssueUnpin(ctx *context.Context) {
	issue, err := issues_model.GetIssueByIndex(ctx, ctx.Repo.Repository.ID, ctx.PathParamInt64("index"))
	if err != nil {
		ctx.ServerError("GetIssueByIndex", err)
		return
	}

	// If we don't do this, it will crash when trying to add the pin event to the comment history
	err = issue.LoadRepo(ctx)
	if err != nil {
		ctx.ServerError("LoadRepo", err)
		return
	}

	err = issues_model.UnpinIssue(ctx, issue, ctx.Doer)
	if err != nil {
		ctx.ServerError("UnpinIssue", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// IssuePinMove moves a pinned Issue
func IssuePinMove(ctx *context.Context) {
	if ctx.Doer == nil {
		ctx.JSON(http.StatusForbidden, "Only signed in users are allowed to perform this action.")
		return
	}

	type movePinIssueForm struct {
		ID       int64 `json:"id"`
		Position int   `json:"position"`
	}

	form := &movePinIssueForm{}
	if err := json.NewDecoder(ctx.Req.Body).Decode(&form); err != nil {
		ctx.ServerError("Decode", err)
		return
	}

	issue, err := issues_model.GetIssueByID(ctx, form.ID)
	if err != nil {
		ctx.ServerError("GetIssueByID", err)
		return
	}

	if issue.RepoID != ctx.Repo.Repository.ID {
		ctx.Status(http.StatusNotFound)
		log.Error("Issue does not belong to this repository")
		return
	}

	err = issues_model.MovePin(ctx, issue, form.Position)
	if err != nil {
		ctx.ServerError("MovePin", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
