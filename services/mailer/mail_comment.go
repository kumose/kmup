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

package mailer

import (
	"context"

	activities_model "github.com/kumose/kmup/models/activities"
	issues_model "github.com/kumose/kmup/models/issues"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

// MailParticipantsComment sends new comment emails to repository watchers and mentioned people.
func MailParticipantsComment(ctx context.Context, c *issues_model.Comment, opType activities_model.ActionType, issue *issues_model.Issue, mentions []*user_model.User) error {
	if setting.MailService == nil {
		// No mail service configured
		return nil
	}

	content := c.Content
	if c.Type == issues_model.CommentTypePullRequestPush {
		content = ""
	}
	if err := mailIssueCommentToParticipants(ctx,
		&mailComment{
			Issue:      issue,
			Doer:       c.Poster,
			ActionType: opType,
			Content:    content,
			Comment:    c,
		}, mentions); err != nil {
		log.Error("mailIssueCommentToParticipants: %v", err)
	}
	return nil
}

// MailMentionsComment sends email to users mentioned in a code comment
func MailMentionsComment(ctx context.Context, pr *issues_model.PullRequest, c *issues_model.Comment, mentions []*user_model.User) (err error) {
	if setting.MailService == nil {
		// No mail service configured
		return nil
	}

	visited := make(container.Set[int64], len(mentions)+1)
	visited.Add(c.Poster.ID)
	if err = mailIssueCommentBatch(ctx,
		&mailComment{
			Issue:      pr.Issue,
			Doer:       c.Poster,
			ActionType: activities_model.ActionCommentPull,
			Content:    c.Content,
			Comment:    c,
		}, mentions, visited, true); err != nil {
		log.Error("mailIssueCommentBatch: %v", err)
	}
	return nil
}
