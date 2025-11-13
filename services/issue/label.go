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

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	access_model "github.com/kumose/kmup/models/perm/access"
	user_model "github.com/kumose/kmup/models/user"
	notify_service "github.com/kumose/kmup/services/notify"
)

// ClearLabels clears all of an issue's labels
func ClearLabels(ctx context.Context, issue *issues_model.Issue, doer *user_model.User) error {
	if err := issues_model.ClearIssueLabels(ctx, issue, doer); err != nil {
		return err
	}

	notify_service.IssueClearLabels(ctx, doer, issue)

	return nil
}

// AddLabel adds a new label to the issue.
func AddLabel(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, label *issues_model.Label) error {
	if err := issues_model.NewIssueLabel(ctx, issue, label, doer); err != nil {
		return err
	}

	notify_service.IssueChangeLabels(ctx, doer, issue, []*issues_model.Label{label}, nil)
	return nil
}

// AddLabels adds a list of new labels to the issue.
func AddLabels(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, labels []*issues_model.Label) error {
	if err := issues_model.NewIssueLabels(ctx, issue, labels, doer); err != nil {
		return err
	}

	notify_service.IssueChangeLabels(ctx, doer, issue, labels, nil)
	return nil
}

// RemoveLabel removes a label from issue by given ID.
func RemoveLabel(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, label *issues_model.Label) error {
	if err := db.WithTx(ctx, func(ctx context.Context) error {
		if err := issue.LoadRepo(ctx); err != nil {
			return err
		}

		perm, err := access_model.GetUserRepoPermission(ctx, issue.Repo, doer)
		if err != nil {
			return err
		}
		if !perm.CanWriteIssuesOrPulls(issue.IsPull) {
			if label.OrgID > 0 {
				return issues_model.ErrOrgLabelNotExist{}
			}
			return issues_model.ErrRepoLabelNotExist{}
		}

		return issues_model.DeleteIssueLabel(ctx, issue, label, doer)
	}); err != nil {
		return err
	}

	notify_service.IssueChangeLabels(ctx, doer, issue, nil, []*issues_model.Label{label})
	return nil
}

// ReplaceLabels removes all current labels and add new labels to the issue.
func ReplaceLabels(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, labels []*issues_model.Label) error {
	old, err := issues_model.GetLabelsByIssueID(ctx, issue.ID)
	if err != nil {
		return err
	}

	if err := issues_model.ReplaceIssueLabels(ctx, issue, labels, doer); err != nil {
		return err
	}

	notify_service.IssueChangeLabels(ctx, doer, issue, labels, old)
	return nil
}
