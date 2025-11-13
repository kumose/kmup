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
	"testing"

	activities_model "github.com/kumose/kmup/models/activities"
	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/repository"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestUpdateIssuesCommit(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	pushCommits := []*repository.PushCommit{
		{
			Sha1:           "abcdef1",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user4@example.com",
			AuthorName:     "User Four",
			Message:        "start working on #FST-1, #1",
		},
		{
			Sha1:           "abcdef2",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user2@example.com",
			AuthorName:     "User Two",
			Message:        "a plain message",
		},
		{
			Sha1:           "abcdef2",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user2@example.com",
			AuthorName:     "User Two",
			Message:        "close #2",
		},
	}

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	repo.Owner = user

	commentBean := &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef1",
		PosterID:  user.ID,
		IssueID:   1,
	}
	issueBean := &issues_model.Issue{RepoID: repo.ID, Index: 4}

	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, &issues_model.Issue{RepoID: repo.ID, Index: 2}, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, repo.DefaultBranch))
	unittest.AssertExistsAndLoadBean(t, commentBean)
	unittest.AssertExistsAndLoadBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})

	// Test that push to a non-default branch closes no issue.
	pushCommits = []*repository.PushCommit{
		{
			Sha1:           "abcdef1",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user4@example.com",
			AuthorName:     "User Four",
			Message:        "close #1",
		},
	}
	repo = unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
	commentBean = &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef1",
		PosterID:  user.ID,
		IssueID:   6,
	}
	issueBean = &issues_model.Issue{RepoID: repo.ID, Index: 1}

	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, &issues_model.Issue{RepoID: repo.ID, Index: 1}, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, "non-existing-branch"))
	unittest.AssertExistsAndLoadBean(t, commentBean)
	unittest.AssertNotExistsBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})

	pushCommits = []*repository.PushCommit{
		{
			Sha1:           "abcdef3",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user2@example.com",
			AuthorName:     "User Two",
			Message:        "close " + setting.AppURL + repo.FullName() + "/pulls/1",
		},
	}
	repo = unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
	commentBean = &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef3",
		PosterID:  user.ID,
		IssueID:   6,
	}
	issueBean = &issues_model.Issue{RepoID: repo.ID, Index: 1}

	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, &issues_model.Issue{RepoID: repo.ID, Index: 1}, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, repo.DefaultBranch))
	unittest.AssertExistsAndLoadBean(t, commentBean)
	unittest.AssertExistsAndLoadBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})
}

func TestUpdateIssuesCommit_Colon(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	pushCommits := []*repository.PushCommit{
		{
			Sha1:           "abcdef2",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user2@example.com",
			AuthorName:     "User Two",
			Message:        "close: #2",
		},
	}

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	repo.Owner = user

	issueBean := &issues_model.Issue{RepoID: repo.ID, Index: 4}

	unittest.AssertNotExistsBean(t, &issues_model.Issue{RepoID: repo.ID, Index: 2}, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, repo.DefaultBranch))
	unittest.AssertExistsAndLoadBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})
}

func TestUpdateIssuesCommit_Issue5957(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	// Test that push to a non-default branch closes an issue.
	pushCommits := []*repository.PushCommit{
		{
			Sha1:           "abcdef1",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user4@example.com",
			AuthorName:     "User Four",
			Message:        "close #2",
		},
	}

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	commentBean := &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef1",
		PosterID:  user.ID,
		IssueID:   7,
	}

	issueBean := &issues_model.Issue{RepoID: repo.ID, Index: 2, ID: 7}

	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, issueBean, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, "non-existing-branch"))
	unittest.AssertExistsAndLoadBean(t, commentBean)
	unittest.AssertExistsAndLoadBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})
}

func TestUpdateIssuesCommit_AnotherRepo(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	// Test that a push to default branch closes issue in another repo
	// If the user also has push permissions to that repo
	pushCommits := []*repository.PushCommit{
		{
			Sha1:           "abcdef1",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user2@example.com",
			AuthorName:     "User Two",
			Message:        "close user2/repo1#1",
		},
	}

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	commentBean := &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef1",
		PosterID:  user.ID,
		IssueID:   1,
	}

	issueBean := &issues_model.Issue{RepoID: 1, Index: 1, ID: 1}

	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, issueBean, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, repo.DefaultBranch))
	unittest.AssertExistsAndLoadBean(t, commentBean)
	unittest.AssertExistsAndLoadBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})
}

func TestUpdateIssuesCommit_AnotherRepo_FullAddress(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	// Test that a push to default branch closes issue in another repo
	// If the user also has push permissions to that repo
	pushCommits := []*repository.PushCommit{
		{
			Sha1:           "abcdef1",
			CommitterEmail: "user2@example.com",
			CommitterName:  "User Two",
			AuthorEmail:    "user2@example.com",
			AuthorName:     "User Two",
			Message:        "close " + setting.AppURL + "user2/repo1/issues/1",
		},
	}

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	commentBean := &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef1",
		PosterID:  user.ID,
		IssueID:   1,
	}

	issueBean := &issues_model.Issue{RepoID: 1, Index: 1, ID: 1}

	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, issueBean, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, repo.DefaultBranch))
	unittest.AssertExistsAndLoadBean(t, commentBean)
	unittest.AssertExistsAndLoadBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})
}

func TestUpdateIssuesCommit_AnotherRepoNoPermission(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 10})

	// Test that a push with close reference *can not* close issue
	// If the committer doesn't have push rights in that repo
	pushCommits := []*repository.PushCommit{
		{
			Sha1:           "abcdef3",
			CommitterEmail: "user10@example.com",
			CommitterName:  "User Ten",
			AuthorEmail:    "user10@example.com",
			AuthorName:     "User Ten",
			Message:        "close org3/repo3#1",
		},
		{
			Sha1:           "abcdef4",
			CommitterEmail: "user10@example.com",
			CommitterName:  "User Ten",
			AuthorEmail:    "user10@example.com",
			AuthorName:     "User Ten",
			Message:        "close " + setting.AppURL + "org3/repo3/issues/1",
		},
	}

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 6})
	commentBean := &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef3",
		PosterID:  user.ID,
		IssueID:   6,
	}
	commentBean2 := &issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		CommitSHA: "abcdef4",
		PosterID:  user.ID,
		IssueID:   6,
	}

	issueBean := &issues_model.Issue{RepoID: 3, Index: 1, ID: 6}

	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, commentBean2)
	unittest.AssertNotExistsBean(t, issueBean, "is_closed=1")
	assert.NoError(t, UpdateIssuesCommit(t.Context(), user, repo, pushCommits, repo.DefaultBranch))
	unittest.AssertNotExistsBean(t, commentBean)
	unittest.AssertNotExistsBean(t, commentBean2)
	unittest.AssertNotExistsBean(t, issueBean, "is_closed=1")
	unittest.CheckConsistencyFor(t, &activities_model.Action{})
}
