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

package pull_test

import (
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	pull_service "github.com/kumose/kmup/services/pull"

	"github.com/stretchr/testify/assert"
)

func TestDismissReview(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	pull := unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{})
	assert.NoError(t, pull.LoadIssue(t.Context()))
	issue := pull.Issue
	assert.NoError(t, issue.LoadRepo(t.Context()))
	reviewer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	review, err := issues_model.CreateReview(t.Context(), issues_model.CreateReviewOptions{
		Issue:    issue,
		Reviewer: reviewer,
		Type:     issues_model.ReviewTypeReject,
	})

	assert.NoError(t, err)
	issue.IsClosed = true
	pull.HasMerged = false
	assert.NoError(t, issues_model.UpdateIssueCols(t.Context(), issue, "is_closed"))
	assert.NoError(t, pull.UpdateCols(t.Context(), "has_merged"))
	_, err = pull_service.DismissReview(t.Context(), review.ID, issue.RepoID, "", &user_model.User{}, false, false)
	assert.Error(t, err)
	assert.True(t, pull_service.IsErrDismissRequestOnClosedPR(err))

	pull.HasMerged = true
	pull.Issue.IsClosed = false
	assert.NoError(t, issues_model.UpdateIssueCols(t.Context(), issue, "is_closed"))
	assert.NoError(t, pull.UpdateCols(t.Context(), "has_merged"))
	_, err = pull_service.DismissReview(t.Context(), review.ID, issue.RepoID, "", &user_model.User{}, false, false)
	assert.Error(t, err)
	assert.True(t, pull_service.IsErrDismissRequestOnClosedPR(err))
}
