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
	"testing"

	webhook_module "github.com/kumose/kmup/modules/webhook"

	"github.com/stretchr/testify/assert"
)

func TestCanGithubEventMatch(t *testing.T) {
	testCases := []struct {
		desc           string
		eventName      string
		triggeredEvent webhook_module.HookEventType
		expected       bool
	}{
		// registry_package event
		{
			"registry_package matches",
			GithubEventRegistryPackage,
			webhook_module.HookEventPackage,
			true,
		},
		{
			"registry_package cannot match",
			GithubEventRegistryPackage,
			webhook_module.HookEventPush,
			false,
		},
		// issues event
		{
			"issue matches",
			GithubEventIssues,
			webhook_module.HookEventIssueLabel,
			true,
		},
		{
			"issue cannot match",
			GithubEventIssues,
			webhook_module.HookEventIssueComment,
			false,
		},
		// issue_comment event
		{
			"issue_comment matches",
			GithubEventIssueComment,
			webhook_module.HookEventIssueComment,
			true,
		},
		{
			"issue_comment cannot match",
			GithubEventIssueComment,
			webhook_module.HookEventIssues,
			false,
		},
		// pull_request event
		{
			"pull_request matches",
			GithubEventPullRequest,
			webhook_module.HookEventPullRequestSync,
			true,
		},
		{
			"pull_request cannot match",
			GithubEventPullRequest,
			webhook_module.HookEventPullRequestComment,
			false,
		},
		// pull_request_target event
		{
			"pull_request_target matches",
			GithubEventPullRequest,
			webhook_module.HookEventPullRequest,
			true,
		},
		{
			"pull_request_target cannot match",
			GithubEventPullRequest,
			webhook_module.HookEventPullRequestComment,
			false,
		},
		// pull_request_review event
		{
			"pull_request_review matches",
			GithubEventPullRequestReview,
			webhook_module.HookEventPullRequestReviewComment,
			true,
		},
		{
			"pull_request_review cannot match",
			GithubEventPullRequestReview,
			webhook_module.HookEventPullRequestComment,
			false,
		},
		// other events
		{
			"create event",
			GithubEventCreate,
			webhook_module.HookEventCreate,
			true,
		},
		{
			"create pull request comment",
			GithubEventIssueComment,
			webhook_module.HookEventPullRequestComment,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equalf(t, tc.expected, canGithubEventMatch(tc.eventName, tc.triggeredEvent), "canGithubEventMatch(%v, %v)", tc.eventName, tc.triggeredEvent)
		})
	}
}
