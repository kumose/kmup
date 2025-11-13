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
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestIssueList_LoadRepositories(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	issueList := issues_model.IssueList{
		unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 4}),
	}

	repos, err := issueList.LoadRepositories(t.Context())
	assert.NoError(t, err)
	assert.Len(t, repos, 2)
	for _, issue := range issueList {
		assert.Equal(t, issue.RepoID, issue.Repo.ID)
	}
}

func TestIssueList_LoadAttributes(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	setting.Service.EnableTimetracking = true
	issueList := issues_model.IssueList{
		unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 4}),
	}

	assert.NoError(t, issueList.LoadAttributes(t.Context()))
	for _, issue := range issueList {
		assert.Equal(t, issue.RepoID, issue.Repo.ID)
		for _, label := range issue.Labels {
			assert.Equal(t, issue.RepoID, label.RepoID)
			unittest.AssertExistsAndLoadBean(t, &issues_model.IssueLabel{IssueID: issue.ID, LabelID: label.ID})
		}
		if issue.PosterID > 0 {
			assert.Equal(t, issue.PosterID, issue.Poster.ID)
		}
		if issue.AssigneeID > 0 {
			assert.Equal(t, issue.AssigneeID, issue.Assignee.ID)
		}
		if issue.MilestoneID > 0 {
			assert.Equal(t, issue.MilestoneID, issue.Milestone.ID)
		}
		if issue.IsPull {
			assert.Equal(t, issue.ID, issue.PullRequest.IssueID)
		}
		for _, attachment := range issue.Attachments {
			assert.Equal(t, issue.ID, attachment.IssueID)
		}
		for _, comment := range issue.Comments {
			assert.Equal(t, issue.ID, comment.IssueID)
		}
		if issue.ID == int64(1) {
			assert.Equal(t, int64(400), issue.TotalTrackedTime)
			assert.NotNil(t, issue.Project)
			assert.Equal(t, int64(1), issue.Project.ID)
		} else {
			assert.Nil(t, issue.Project)
		}
	}
}
