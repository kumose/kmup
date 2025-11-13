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
	access_model "github.com/kumose/kmup/models/perm/access"
	user_model "github.com/kumose/kmup/models/user"
	notify_service "github.com/kumose/kmup/services/notify"
)

// ChangeContent changes issue content, as the given user.
func ChangeContent(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, content string, contentVersion int) error {
	if err := issue.LoadRepo(ctx); err != nil {
		return err
	}

	if user_model.IsUserBlockedBy(ctx, doer, issue.PosterID, issue.Repo.OwnerID) {
		if isAdmin, _ := access_model.IsUserRepoAdmin(ctx, issue.Repo, doer); !isAdmin {
			return user_model.ErrBlockedUser
		}
	}

	oldContent := issue.Content

	if err := issues_model.ChangeIssueContent(ctx, issue, doer, content, contentVersion); err != nil {
		return err
	}

	notify_service.IssueChangeContent(ctx, doer, issue, oldContent)

	return nil
}
