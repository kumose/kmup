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
	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/eventsource"
	"github.com/kumose/kmup/services/context"
)

// IssueStartStopwatch creates a stopwatch for the given issue.
func IssueStartStopwatch(c *context.Context) {
	issue := GetActionIssue(c)
	if c.Written() {
		return
	}

	if !c.Repo.CanUseTimetracker(c, issue, c.Doer) {
		c.NotFound(nil)
		return
	}

	if ok, err := issues_model.CreateIssueStopwatch(c, c.Doer, issue); err != nil {
		c.ServerError("CreateIssueStopwatch", err)
		return
	} else if !ok {
		c.Flash.Warning(c.Tr("repo.issues.stopwatch_already_created"))
	} else {
		c.Flash.Success(c.Tr("repo.issues.tracker_auto_close"))
	}
	c.JSONRedirect("")
}

// IssueStopStopwatch stops a stopwatch for the given issue.
func IssueStopStopwatch(c *context.Context) {
	issue := GetActionIssue(c)
	if c.Written() {
		return
	}

	if !c.Repo.CanUseTimetracker(c, issue, c.Doer) {
		c.NotFound(nil)
		return
	}

	if ok, err := issues_model.FinishIssueStopwatch(c, c.Doer, issue); err != nil {
		c.ServerError("FinishIssueStopwatch", err)
		return
	} else if !ok {
		c.Flash.Warning(c.Tr("repo.issues.stopwatch_already_stopped"))
	}
	c.JSONRedirect("")
}

// CancelStopwatch cancel the stopwatch
func CancelStopwatch(c *context.Context) {
	issue := GetActionIssue(c)
	if c.Written() {
		return
	}
	if !c.Repo.CanUseTimetracker(c, issue, c.Doer) {
		c.NotFound(nil)
		return
	}

	if _, err := issues_model.CancelStopwatch(c, c.Doer, issue); err != nil {
		c.ServerError("CancelStopwatch", err)
		return
	}

	stopwatches, err := issues_model.GetUserStopwatches(c, c.Doer.ID, db.ListOptions{})
	if err != nil {
		c.ServerError("GetUserStopwatches", err)
		return
	}
	if len(stopwatches) == 0 {
		eventsource.GetManager().SendMessage(c.Doer.ID, &eventsource.Event{
			Name: "stopwatches",
			Data: "{}",
		})
	}

	c.JSONRedirect("")
}
