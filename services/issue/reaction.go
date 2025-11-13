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

package issue

import (
	"context"

	issues_model "github.com/kumose/kmup/models/issues"
	user_model "github.com/kumose/kmup/models/user"
)

// CreateIssueReaction creates a reaction on an issue.
func CreateIssueReaction(ctx context.Context, doer *user_model.User, issue *issues_model.Issue, content string) (*issues_model.Reaction, error) {
	if err := issue.LoadRepo(ctx); err != nil {
		return nil, err
	}

	if user_model.IsUserBlockedBy(ctx, doer, issue.PosterID, issue.Repo.OwnerID) {
		return nil, user_model.ErrBlockedUser
	}

	return issues_model.CreateReaction(ctx, &issues_model.ReactionOptions{
		Type:    content,
		DoerID:  doer.ID,
		IssueID: issue.ID,
	})
}

// CreateCommentReaction creates a reaction on a comment.
func CreateCommentReaction(ctx context.Context, doer *user_model.User, comment *issues_model.Comment, content string) (*issues_model.Reaction, error) {
	if err := comment.LoadIssue(ctx); err != nil {
		return nil, err
	}

	if err := comment.Issue.LoadRepo(ctx); err != nil {
		return nil, err
	}

	if user_model.IsUserBlockedBy(ctx, doer, comment.Issue.PosterID, comment.Issue.Repo.OwnerID, comment.PosterID) {
		return nil, user_model.ErrBlockedUser
	}

	return issues_model.CreateReaction(ctx, &issues_model.ReactionOptions{
		Type:      content,
		DoerID:    doer.ID,
		IssueID:   comment.Issue.ID,
		CommentID: comment.ID,
	})
}
