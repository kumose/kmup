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

package integration

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	auth_model "github.com/kumose/kmup/models/auth"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/gitrepo"
	pull_service "github.com/kumose/kmup/services/pull"
	repo_service "github.com/kumose/kmup/services/repository"
	files_service "github.com/kumose/kmup/services/repository/files"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIPullUpdate(t *testing.T) {
	onKmupRun(t, func(t *testing.T, kmupURL *url.URL) {
		// Create PR to test
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		org26 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 26})
		pr := createOutdatedPR(t, user, org26)
		require.NoError(t, pr.LoadBaseRepo(t.Context()))
		require.NoError(t, pr.LoadIssue(t.Context()))

		// Test GetDiverging
		diffCount, err := gitrepo.GetDivergingCommits(t.Context(), pr.BaseRepo, pr.BaseBranch, pr.GetGitHeadRefName())
		require.NoError(t, err)
		assert.Equal(t, 1, diffCount.Behind)
		assert.Equal(t, 1, diffCount.Ahead)
		assert.Equal(t, diffCount.Behind, pr.CommitsBehind)
		assert.Equal(t, diffCount.Ahead, pr.CommitsAhead)

		session := loginUser(t, "user2")
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
		req := NewRequestf(t, "POST", "/api/v1/repos/%s/%s/pulls/%d/update", pr.BaseRepo.OwnerName, pr.BaseRepo.Name, pr.Issue.Index).
			AddTokenAuth(token)
		session.MakeRequest(t, req, http.StatusOK)

		// Test GetDiverging after update
		diffCount, err = gitrepo.GetDivergingCommits(t.Context(), pr.BaseRepo, pr.BaseBranch, pr.GetGitHeadRefName())
		require.NoError(t, err)
		assert.Equal(t, 0, diffCount.Behind)
		assert.Equal(t, 2, diffCount.Ahead)
		assert.Eventually(t, func() bool {
			pr := unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: pr.ID})
			return diffCount.Behind == pr.CommitsBehind && diffCount.Ahead == pr.CommitsAhead
		}, 5*time.Second, 20*time.Millisecond)
	})
}

func TestAPIPullUpdateByRebase(t *testing.T) {
	onKmupRun(t, func(t *testing.T, kmupURL *url.URL) {
		// Create PR to test
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		org26 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 26})
		pr := createOutdatedPR(t, user, org26)
		assert.NoError(t, pr.LoadBaseRepo(t.Context()))

		// Test GetDiverging
		diffCount, err := gitrepo.GetDivergingCommits(t.Context(), pr.BaseRepo, pr.BaseBranch, pr.GetGitHeadRefName())
		assert.NoError(t, err)
		assert.Equal(t, 1, diffCount.Behind)
		assert.Equal(t, 1, diffCount.Ahead)
		assert.NoError(t, pr.LoadIssue(t.Context()))

		session := loginUser(t, "user2")
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)
		req := NewRequestf(t, "POST", "/api/v1/repos/%s/%s/pulls/%d/update?style=rebase", pr.BaseRepo.OwnerName, pr.BaseRepo.Name, pr.Issue.Index).
			AddTokenAuth(token)
		session.MakeRequest(t, req, http.StatusOK)

		// Test GetDiverging after update
		diffCount, err = gitrepo.GetDivergingCommits(t.Context(), pr.BaseRepo, pr.BaseBranch, pr.GetGitHeadRefName())
		assert.NoError(t, err)
		assert.Equal(t, 0, diffCount.Behind)
		assert.Equal(t, 1, diffCount.Ahead)
	})
}

func createOutdatedPR(t *testing.T, actor, forkOrg *user_model.User) *issues_model.PullRequest {
	baseRepo, err := repo_service.CreateRepository(t.Context(), actor, actor, repo_service.CreateRepoOptions{
		Name:        "repo-pr-update",
		Description: "repo-tmp-pr-update description",
		AutoInit:    true,
		Gitignores:  "C,C++",
		License:     "MIT",
		Readme:      "Default",
		IsPrivate:   false,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, baseRepo)

	headRepo, err := repo_service.ForkRepository(t.Context(), actor, forkOrg, repo_service.ForkRepoOptions{
		BaseRepo:    baseRepo,
		Name:        "repo-pr-update",
		Description: "desc",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, headRepo)

	// create a commit on base Repo
	_, err = files_service.ChangeRepoFiles(t.Context(), baseRepo, actor, &files_service.ChangeRepoFilesOptions{
		Files: []*files_service.ChangeRepoFile{
			{
				Operation:     "create",
				TreePath:      "File_A",
				ContentReader: strings.NewReader("File A"),
			},
		},
		Message:   "Add File A",
		OldBranch: "master",
		NewBranch: "master",
		Author: &files_service.IdentityOptions{
			GitUserName:  actor.Name,
			GitUserEmail: actor.Email,
		},
		Committer: &files_service.IdentityOptions{
			GitUserName:  actor.Name,
			GitUserEmail: actor.Email,
		},
		Dates: &files_service.CommitDateOptions{
			Author:    time.Now(),
			Committer: time.Now(),
		},
	})
	assert.NoError(t, err)

	// create a commit on head Repo
	_, err = files_service.ChangeRepoFiles(t.Context(), headRepo, actor, &files_service.ChangeRepoFilesOptions{
		Files: []*files_service.ChangeRepoFile{
			{
				Operation:     "create",
				TreePath:      "File_B",
				ContentReader: strings.NewReader("File B"),
			},
		},
		Message:   "Add File on PR branch",
		OldBranch: "master",
		NewBranch: "newBranch",
		Author: &files_service.IdentityOptions{
			GitUserName:  actor.Name,
			GitUserEmail: actor.Email,
		},
		Committer: &files_service.IdentityOptions{
			GitUserName:  actor.Name,
			GitUserEmail: actor.Email,
		},
		Dates: &files_service.CommitDateOptions{
			Author:    time.Now(),
			Committer: time.Now(),
		},
	})
	assert.NoError(t, err)

	// create Pull
	pullIssue := &issues_model.Issue{
		RepoID:   baseRepo.ID,
		Title:    "Test Pull -to-update-",
		PosterID: actor.ID,
		Poster:   actor,
		IsPull:   true,
	}
	pullRequest := &issues_model.PullRequest{
		HeadRepoID: headRepo.ID,
		BaseRepoID: baseRepo.ID,
		HeadBranch: "newBranch",
		BaseBranch: "master",
		HeadRepo:   headRepo,
		BaseRepo:   baseRepo,
		Type:       issues_model.PullRequestKmup,
	}
	prOpts := &pull_service.NewPullRequestOptions{Repo: baseRepo, Issue: pullIssue, PullRequest: pullRequest}
	err = pull_service.NewPullRequest(t.Context(), prOpts)
	assert.NoError(t, err)

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{Title: "Test Pull -to-update-"})
	assert.NoError(t, issue.LoadPullRequest(t.Context()))

	return issue.PullRequest
}
