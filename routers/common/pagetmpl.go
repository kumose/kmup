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

package common

import (
	goctx "context"
	"errors"
	"sync"

	activities_model "github.com/kumose/kmup/models/activities"
	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/context"
)

// StopwatchTmplInfo is a view on a stopwatch specifically for template rendering
type StopwatchTmplInfo struct {
	IssueLink  string
	RepoSlug   string
	IssueIndex int64
	Seconds    int64
}

func getActiveStopwatch(ctx *context.Context) *StopwatchTmplInfo {
	if ctx.Doer == nil {
		return nil
	}

	_, sw, issue, err := issues_model.HasUserStopwatch(ctx, ctx.Doer.ID)
	if err != nil {
		if !errors.Is(err, goctx.Canceled) {
			log.Error("Unable to HasUserStopwatch for user:%-v: %v", ctx.Doer, err)
		}
		return nil
	}

	if sw == nil || sw.ID == 0 {
		return nil
	}

	return &StopwatchTmplInfo{
		issue.Link(),
		issue.Repo.FullName(),
		issue.Index,
		sw.Seconds() + 1, // ensure time is never zero in ui
	}
}

func notificationUnreadCount(ctx *context.Context) int64 {
	if ctx.Doer == nil {
		return 0
	}
	count, err := db.Count[activities_model.Notification](ctx, activities_model.FindNotificationOptions{
		UserID: ctx.Doer.ID,
		Status: []activities_model.NotificationStatus{activities_model.NotificationStatusUnread},
	})
	if err != nil {
		if !errors.Is(err, goctx.Canceled) {
			log.Error("Unable to find notification for user:%-v: %v", ctx.Doer, err)
		}
		return 0
	}
	return count
}

type pageGlobalDataType struct {
	IsSigned    bool
	IsSiteAdmin bool

	GetNotificationUnreadCount func() int64
	GetActiveStopwatch         func() *StopwatchTmplInfo
}

func PageGlobalData(ctx *context.Context) {
	var data pageGlobalDataType
	data.IsSigned = ctx.Doer != nil
	data.IsSiteAdmin = ctx.Doer != nil && ctx.Doer.IsAdmin
	data.GetNotificationUnreadCount = sync.OnceValue(func() int64 { return notificationUnreadCount(ctx) })
	data.GetActiveStopwatch = sync.OnceValue(func() *StopwatchTmplInfo { return getActiveStopwatch(ctx) })
	ctx.Data["PageGlobalData"] = data
}
