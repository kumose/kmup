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
	"errors"
	"fmt"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	user_model "github.com/kumose/kmup/models/user"
	notify_service "github.com/kumose/kmup/services/notify"
)

func changeMilestoneAssign(ctx context.Context, doer *user_model.User, issue *issues_model.Issue, oldMilestoneID int64) error {
	// Only check if milestone exists if we don't remove it.
	if issue.MilestoneID > 0 {
		has, err := issues_model.HasMilestoneByRepoID(ctx, issue.RepoID, issue.MilestoneID)
		if err != nil {
			return fmt.Errorf("HasMilestoneByRepoID: %w", err)
		}
		if !has {
			return errors.New("HasMilestoneByRepoID: issue doesn't exist")
		}
	}

	if err := issues_model.UpdateIssueCols(ctx, issue, "milestone_id"); err != nil {
		return err
	}

	if oldMilestoneID > 0 {
		if err := issues_model.UpdateMilestoneCounters(ctx, oldMilestoneID); err != nil {
			return err
		}
	}

	if issue.MilestoneID > 0 {
		if err := issues_model.UpdateMilestoneCounters(ctx, issue.MilestoneID); err != nil {
			return err
		}
	}

	if oldMilestoneID > 0 || issue.MilestoneID > 0 {
		if err := issue.LoadRepo(ctx); err != nil {
			return err
		}

		opts := &issues_model.CreateCommentOptions{
			Type:           issues_model.CommentTypeMilestone,
			Doer:           doer,
			Repo:           issue.Repo,
			Issue:          issue,
			OldMilestoneID: oldMilestoneID,
			MilestoneID:    issue.MilestoneID,
		}
		if _, err := issues_model.CreateComment(ctx, opts); err != nil {
			return err
		}
	}

	if issue.MilestoneID == 0 {
		issue.Milestone = nil
	}

	return nil
}

// ChangeMilestoneAssign changes assignment of milestone for issue.
func ChangeMilestoneAssign(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, oldMilestoneID int64) (err error) {
	if err := db.WithTx(ctx, func(dbCtx context.Context) error {
		return changeMilestoneAssign(dbCtx, doer, issue, oldMilestoneID)
	}); err != nil {
		return err
	}

	notify_service.IssueChangeMilestone(ctx, doer, issue, oldMilestoneID)
	return nil
}
