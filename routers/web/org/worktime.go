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

package org

import (
	"net/http"
	"time"

	"github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/modules/templates"
	shared_user "github.com/kumose/kmup/routers/web/shared/user"
	"github.com/kumose/kmup/services/context"
)

const tplByRepos templates.TplName = "org/worktime"

// parseOrgTimes contains functionality that is required in all these functions,
// like parsing the date from the request, setting default dates, etc.
func parseOrgTimes(ctx *context.Context) (unixFrom, unixTo int64) {
	rangeFrom := ctx.FormString("from")
	rangeTo := ctx.FormString("to")
	if rangeFrom == "" {
		rangeFrom = time.Now().Format("2006-01") + "-01" // defaults to start of current month
	}
	if rangeTo == "" {
		rangeTo = time.Now().Format("2006-01-02") // defaults to today
	}

	ctx.Data["RangeFrom"] = rangeFrom
	ctx.Data["RangeTo"] = rangeTo

	timeFrom, err := time.Parse("2006-01-02", rangeFrom)
	if err != nil {
		ctx.ServerError("time.Parse", err)
	}
	timeTo, err := time.Parse("2006-01-02", rangeTo)
	if err != nil {
		ctx.ServerError("time.Parse", err)
	}
	unixFrom = timeFrom.Unix()
	unixTo = timeTo.Add(1440*time.Minute - 1*time.Second).Unix() // humans expect that we include the ending day too
	return unixFrom, unixTo
}

func Worktime(ctx *context.Context) {
	ctx.Data["PageIsOrgTimes"] = true

	unixFrom, unixTo := parseOrgTimes(ctx)
	if ctx.Written() {
		return
	}

	worktimeBy := ctx.FormString("by")
	ctx.Data["WorktimeBy"] = worktimeBy

	var worktimeSumResult any
	var err error
	switch worktimeBy {
	case "milestones":
		worktimeSumResult, err = organization.GetWorktimeByMilestones(ctx, ctx.Org.Organization, unixFrom, unixTo)
		ctx.Data["WorktimeByMilestones"] = true
	case "members":
		worktimeSumResult, err = organization.GetWorktimeByMembers(ctx, ctx.Org.Organization, unixFrom, unixTo)
		ctx.Data["WorktimeByMembers"] = true
	default: /* by repos */
		worktimeSumResult, err = organization.GetWorktimeByRepos(ctx, ctx.Org.Organization, unixFrom, unixTo)
		ctx.Data["WorktimeByRepos"] = true
	}
	if err != nil {
		ctx.ServerError("GetWorktime", err)
		return
	}

	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	ctx.Data["WorktimeSumResult"] = worktimeSumResult
	ctx.HTML(http.StatusOK, tplByRepos)
}
