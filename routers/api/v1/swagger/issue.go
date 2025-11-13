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

package swagger

import (
	api "github.com/kumose/kmup/modules/structs"
)

// Issue
// swagger:response Issue
type swaggerResponseIssue struct {
	// in:body
	Body api.Issue `json:"body"`
}

// IssueList
// swagger:response IssueList
type swaggerResponseIssueList struct {
	// in:body
	Body []api.Issue `json:"body"`
}

// Comment
// swagger:response Comment
type swaggerResponseComment struct {
	// in:body
	Body api.Comment `json:"body"`
}

// CommentList
// swagger:response CommentList
type swaggerResponseCommentList struct {
	// in:body
	Body []api.Comment `json:"body"`
}

// TimelineList
// swagger:response TimelineList
type swaggerResponseTimelineList struct {
	// in:body
	Body []api.TimelineComment `json:"body"`
}

// Label
// swagger:response Label
type swaggerResponseLabel struct {
	// in:body
	Body api.Label `json:"body"`
}

// LabelList
// swagger:response LabelList
type swaggerResponseLabelList struct {
	// in:body
	Body []api.Label `json:"body"`
}

// Milestone
// swagger:response Milestone
type swaggerResponseMilestone struct {
	// in:body
	Body api.Milestone `json:"body"`
}

// MilestoneList
// swagger:response MilestoneList
type swaggerResponseMilestoneList struct {
	// in:body
	Body []api.Milestone `json:"body"`
}

// TrackedTime
// swagger:response TrackedTime
type swaggerResponseTrackedTime struct {
	// in:body
	Body api.TrackedTime `json:"body"`
}

// TrackedTimeList
// swagger:response TrackedTimeList
type swaggerResponseTrackedTimeList struct {
	// in:body
	Body []api.TrackedTime `json:"body"`
}

// IssueDeadline
// swagger:response IssueDeadline
type swaggerIssueDeadline struct {
	// in:body
	Body api.IssueDeadline `json:"body"`
}

// IssueTemplates
// swagger:response IssueTemplates
type swaggerIssueTemplates struct {
	// in:body
	Body []api.IssueTemplate `json:"body"`
}

// StopWatch
// swagger:response StopWatch
type swaggerResponseStopWatch struct {
	// in:body
	Body api.StopWatch `json:"body"`
}

// StopWatchList
// swagger:response StopWatchList
type swaggerResponseStopWatchList struct {
	// in:body
	Body []api.StopWatch `json:"body"`
}

// Reaction
// swagger:response Reaction
type swaggerReaction struct {
	// in:body
	Body api.Reaction `json:"body"`
}

// ReactionList
// swagger:response ReactionList
type swaggerReactionList struct {
	// in:body
	Body []api.Reaction `json:"body"`
}
