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

package issues_test

import (
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestPullRequestList_LoadAttributes(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	prs := issues_model.PullRequestList{
		unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 2}),
	}
	assert.NoError(t, prs.LoadAttributes(t.Context()))
	for _, pr := range prs {
		assert.NotNil(t, pr.Issue)
		assert.Equal(t, pr.IssueID, pr.Issue.ID)
	}

	assert.NoError(t, issues_model.PullRequestList([]*issues_model.PullRequest{}).LoadAttributes(t.Context()))
}

func TestPullRequestList_LoadReviewCommentsCounts(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	prs := issues_model.PullRequestList{
		unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 2}),
	}
	reviewComments, err := prs.LoadReviewCommentsCounts(t.Context())
	assert.NoError(t, err)
	assert.Len(t, reviewComments, 2)
	for _, pr := range prs {
		assert.Equal(t, 1, reviewComments[pr.IssueID])
	}
}

func TestPullRequestList_LoadReviews(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	prs := issues_model.PullRequestList{
		unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 2}),
	}
	reviewList, err := prs.LoadReviews(t.Context())
	assert.NoError(t, err)
	// 1, 7, 8, 9, 10, 22
	assert.Len(t, reviewList, 6)
	assert.EqualValues(t, 1, reviewList[0].ID)
	assert.EqualValues(t, 7, reviewList[1].ID)
	assert.EqualValues(t, 8, reviewList[2].ID)
	assert.EqualValues(t, 9, reviewList[3].ID)
	assert.EqualValues(t, 10, reviewList[4].ID)
	assert.EqualValues(t, 22, reviewList[5].ID)
}
