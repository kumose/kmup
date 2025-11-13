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

package convert

import (
	"context"
	"strings"

	issues_model "github.com/kumose/kmup/models/issues"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
)

// ToPullReview convert a review to api format
func ToPullReview(ctx context.Context, r *issues_model.Review, doer *user_model.User) (*api.PullReview, error) {
	if err := r.LoadAttributes(ctx); err != nil {
		if !user_model.IsErrUserNotExist(err) {
			return nil, err
		}
		r.Reviewer = user_model.NewGhostUser()
	}

	result := &api.PullReview{
		ID:                r.ID,
		Reviewer:          ToUser(ctx, r.Reviewer, doer),
		State:             api.ReviewStateUnknown,
		Body:              r.Content,
		CommitID:          r.CommitID,
		Stale:             r.Stale,
		Official:          r.Official,
		Dismissed:         r.Dismissed,
		CodeCommentsCount: r.GetCodeCommentsCount(ctx),
		Submitted:         r.CreatedUnix.AsTime(),
		Updated:           r.UpdatedUnix.AsTime(),
		HTMLURL:           r.HTMLURL(ctx),
		HTMLPullURL:       r.Issue.HTMLURL(ctx),
	}

	if r.ReviewerTeam != nil {
		var err error
		result.ReviewerTeam, err = ToTeam(ctx, r.ReviewerTeam)
		if err != nil {
			return nil, err
		}
	}

	switch r.Type {
	case issues_model.ReviewTypeApprove:
		result.State = api.ReviewStateApproved
	case issues_model.ReviewTypeReject:
		result.State = api.ReviewStateRequestChanges
	case issues_model.ReviewTypeComment:
		result.State = api.ReviewStateComment
	case issues_model.ReviewTypePending:
		result.State = api.ReviewStatePending
	case issues_model.ReviewTypeRequest:
		result.State = api.ReviewStateRequestReview
	}

	return result, nil
}

// ToPullReviewList convert a list of review to it's api format
func ToPullReviewList(ctx context.Context, rl []*issues_model.Review, doer *user_model.User) ([]*api.PullReview, error) {
	result := make([]*api.PullReview, 0, len(rl))
	for i := range rl {
		// show pending reviews only for the user who created them
		if rl[i].Type == issues_model.ReviewTypePending && (doer == nil || (!doer.IsAdmin && doer.ID != rl[i].ReviewerID)) {
			continue
		}
		r, err := ToPullReview(ctx, rl[i], doer)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

// ToPullReviewCommentList convert the CodeComments of an review to it's api format
func ToPullReviewCommentList(ctx context.Context, review *issues_model.Review, doer *user_model.User) ([]*api.PullReviewComment, error) {
	if err := review.LoadAttributes(ctx); err != nil {
		if !user_model.IsErrUserNotExist(err) {
			return nil, err
		}
		review.Reviewer = user_model.NewGhostUser()
	}

	apiComments := make([]*api.PullReviewComment, 0, len(review.CodeComments))

	for _, lines := range review.CodeComments {
		for _, comments := range lines {
			for _, comment := range comments {
				apiComment := &api.PullReviewComment{
					ID:           comment.ID,
					Body:         comment.Content,
					Poster:       ToUser(ctx, comment.Poster, doer),
					Resolver:     ToUser(ctx, comment.ResolveDoer, doer),
					ReviewID:     review.ID,
					Created:      comment.CreatedUnix.AsTime(),
					Updated:      comment.UpdatedUnix.AsTime(),
					Path:         comment.TreePath,
					CommitID:     comment.CommitSHA,
					OrigCommitID: comment.OldRef,
					DiffHunk:     patch2diff(comment.Patch),
					HTMLURL:      comment.HTMLURL(ctx),
					HTMLPullURL:  review.Issue.HTMLURL(ctx),
				}

				if comment.Line < 0 {
					apiComment.OldLineNum = comment.UnsignedLine()
				} else {
					apiComment.LineNum = comment.UnsignedLine()
				}
				apiComments = append(apiComments, apiComment)
			}
		}
	}
	return apiComments, nil
}

func patch2diff(patch string) string {
	split := strings.Split(patch, "\n@@")
	if len(split) == 2 {
		return "@@" + split[1]
	}
	return ""
}
