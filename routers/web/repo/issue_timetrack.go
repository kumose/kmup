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
	"strings"
	"time"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/forms"
	issue_service "github.com/kumose/kmup/services/issue"
)

// AddTimeManually tracks time manually
func AddTimeManually(c *context.Context) {
	form := web.GetForm(c).(*forms.AddTimeManuallyForm)
	issue := GetActionIssue(c)
	if c.Written() {
		return
	}
	if !c.Repo.CanUseTimetracker(c, issue, c.Doer) {
		c.NotFound(nil)
		return
	}

	if c.HasError() {
		c.JSONError(c.GetErrMsg())
		return
	}

	total := time.Duration(form.Hours)*time.Hour + time.Duration(form.Minutes)*time.Minute

	if total <= 0 {
		c.JSONError(c.Tr("repo.issues.add_time_sum_to_small"))
		return
	}

	if _, err := issues_model.AddTime(c, c.Doer, issue, int64(total.Seconds()), time.Now()); err != nil {
		c.ServerError("AddTime", err)
		return
	}

	c.JSONRedirect("")
}

// DeleteTime deletes tracked time
func DeleteTime(c *context.Context) {
	issue := GetActionIssue(c)
	if c.Written() {
		return
	}
	if !c.Repo.CanUseTimetracker(c, issue, c.Doer) {
		c.NotFound(nil)
		return
	}

	t, err := issues_model.GetTrackedTimeByID(c, c.PathParamInt64("timeid"))
	if err != nil {
		if db.IsErrNotExist(err) {
			c.NotFound(err)
			return
		}
		c.HTTPError(http.StatusInternalServerError, "GetTrackedTimeByID", err.Error())
		return
	}

	// only OP or admin may delete
	if !c.IsSigned || (!c.IsUserSiteAdmin() && c.Doer.ID != t.UserID) {
		c.HTTPError(http.StatusForbidden, "not allowed")
		return
	}

	if err = issues_model.DeleteTime(c, t); err != nil {
		c.ServerError("DeleteTime", err)
		return
	}

	c.Flash.Success(c.Tr("repo.issues.del_time_history", util.SecToHours(t.Time)))
	c.JSONRedirect("")
}

func UpdateIssueTimeEstimate(ctx *context.Context) {
	issue := GetActionIssue(ctx)
	if ctx.Written() {
		return
	}

	if !ctx.IsSigned || (!issue.IsPoster(ctx.Doer.ID) && !ctx.Repo.CanWriteIssuesOrPulls(issue.IsPull)) {
		ctx.HTTPError(http.StatusForbidden)
		return
	}

	timeStr := strings.TrimSpace(ctx.FormString("time_estimate"))

	total, err := util.TimeEstimateParse(timeStr)
	if err != nil {
		ctx.JSONError(ctx.Tr("repo.issues.time_estimate_invalid"))
		return
	}

	// No time changed
	if issue.TimeEstimate == total {
		ctx.JSONRedirect("")
		return
	}

	if err := issue_service.ChangeTimeEstimate(ctx, issue, ctx.Doer, total); err != nil {
		ctx.ServerError("ChangeTimeEstimate", err)
		return
	}

	ctx.JSONRedirect("")
}
