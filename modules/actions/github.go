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

package actions

import (
	webhook_module "github.com/kumose/kmup/modules/webhook"
)

const (
	GithubEventPullRequest              = "pull_request"
	GithubEventPullRequestTarget        = "pull_request_target"
	GithubEventPullRequestReviewComment = "pull_request_review_comment"
	GithubEventPullRequestReview        = "pull_request_review"
	GithubEventRegistryPackage          = "registry_package"
	GithubEventCreate                   = "create"
	GithubEventDelete                   = "delete"
	GithubEventFork                     = "fork"
	GithubEventPush                     = "push"
	GithubEventIssues                   = "issues"
	GithubEventIssueComment             = "issue_comment"
	GithubEventRelease                  = "release"
	GithubEventPullRequestComment       = "pull_request_comment"
	GithubEventGollum                   = "gollum"
	GithubEventSchedule                 = "schedule"
)

// IsDefaultBranchWorkflow returns true if the event only triggers workflows on the default branch
func IsDefaultBranchWorkflow(triggedEvent webhook_module.HookEventType) bool {
	switch triggedEvent {
	case webhook_module.HookEventDelete:
		// GitHub "delete" event
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#delete
		return true
	case webhook_module.HookEventFork:
		// GitHub "fork" event
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#fork
		return true
	case webhook_module.HookEventIssueComment:
		// GitHub "issue_comment" event
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#issue_comment
		return true
	case webhook_module.HookEventPullRequestComment:
		// GitHub "pull_request_comment" event
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request_comment-use-issue_comment
		return true
	case webhook_module.HookEventWiki:
		// GitHub "gollum" event
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#gollum
		return true
	case webhook_module.HookEventSchedule:
		// GitHub "schedule" event
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#schedule
		return true
	case webhook_module.HookEventIssues,
		webhook_module.HookEventIssueAssign,
		webhook_module.HookEventIssueLabel,
		webhook_module.HookEventIssueMilestone:
		// Github "issues" event
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#issues
		return true
	}

	return false
}

// canGithubEventMatch check if the input Github event can match any Kmup event.
func canGithubEventMatch(eventName string, triggedEvent webhook_module.HookEventType) bool {
	switch eventName {
	case GithubEventRegistryPackage:
		return triggedEvent == webhook_module.HookEventPackage

	// See https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#gollum
	case GithubEventGollum:
		return triggedEvent == webhook_module.HookEventWiki

	case GithubEventIssues:
		switch triggedEvent {
		case webhook_module.HookEventIssues,
			webhook_module.HookEventIssueAssign,
			webhook_module.HookEventIssueLabel,
			webhook_module.HookEventIssueMilestone:
			return true

		default:
			return false
		}

	case GithubEventPullRequest, GithubEventPullRequestTarget:
		switch triggedEvent {
		case webhook_module.HookEventPullRequest,
			webhook_module.HookEventPullRequestSync,
			webhook_module.HookEventPullRequestAssign,
			webhook_module.HookEventPullRequestLabel,
			webhook_module.HookEventPullRequestReviewRequest,
			webhook_module.HookEventPullRequestMilestone:
			return true

		default:
			return false
		}

	case GithubEventPullRequestReview:
		switch triggedEvent {
		case webhook_module.HookEventPullRequestReviewApproved,
			webhook_module.HookEventPullRequestReviewComment,
			webhook_module.HookEventPullRequestReviewRejected:
			return true

		default:
			return false
		}

	case GithubEventSchedule:
		return triggedEvent == webhook_module.HookEventSchedule

	case GithubEventIssueComment:
		// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request_comment-use-issue_comment
		return triggedEvent == webhook_module.HookEventIssueComment ||
			triggedEvent == webhook_module.HookEventPullRequestComment

	default:
		return eventName == string(triggedEvent)
	}
}
