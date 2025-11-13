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
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	notify_service "github.com/kumose/kmup/services/notify"
)

// CloseIssue close an issue.
func CloseIssue(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, commitID string) error {
	var comment *issues_model.Comment
	if err := db.WithTx(ctx, func(ctx context.Context) error {
		var err error
		comment, err = issues_model.CloseIssue(ctx, issue, doer)
		if err != nil {
			if issues_model.IsErrDependenciesLeft(err) {
				if _, err := issues_model.FinishIssueStopwatch(ctx, doer, issue); err != nil {
					log.Error("Unable to stop stopwatch for issue[%d]#%d: %v", issue.ID, issue.Index, err)
				}
			}
			return err
		}

		_, err = issues_model.FinishIssueStopwatch(ctx, doer, issue)
		return err
	}); err != nil {
		return err
	}

	notify_service.IssueChangeStatus(ctx, doer, commitID, issue, comment, true)

	return nil
}

// ReopenIssue reopen an issue.
// FIXME: If some issues dependent this one are closed, should we also reopen them?
func ReopenIssue(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, commitID string) error {
	comment, err := issues_model.ReopenIssue(ctx, issue, doer)
	if err != nil {
		return err
	}

	notify_service.IssueChangeStatus(ctx, doer, commitID, issue, comment, false)

	return nil
}
