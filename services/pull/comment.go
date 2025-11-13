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

package pull

import (
	"context"

	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/json"
)

// getCommitIDsFromRepo get commit IDs from repo in between oldCommitID and newCommitID
// Commit on baseBranch will skip
func getCommitIDsFromRepo(ctx context.Context, repo *repo_model.Repository, oldCommitID, newCommitID, baseBranch string) (commitIDs []string, err error) {
	gitRepo, closer, err := gitrepo.RepositoryFromContextOrOpen(ctx, repo)
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	oldCommit, err := gitRepo.GetCommit(oldCommitID)
	if err != nil {
		return nil, err
	}

	newCommit, err := gitRepo.GetCommit(newCommitID)
	if err != nil {
		return nil, err
	}

	// Find commits between new and old commit excluding base branch commits
	commits, err := gitRepo.CommitsBetweenNotBase(newCommit, oldCommit, baseBranch)
	if err != nil {
		return nil, err
	}

	commitIDs = make([]string, 0, len(commits))
	for i := len(commits) - 1; i >= 0; i-- {
		commitIDs = append(commitIDs, commits[i].ID.String())
	}

	return commitIDs, err
}

// CreatePushPullComment create push code to pull base comment
func CreatePushPullComment(ctx context.Context, pusher *user_model.User, pr *issues_model.PullRequest, oldCommitID, newCommitID string, isForcePush bool) (comment *issues_model.Comment, err error) {
	if pr.HasMerged || oldCommitID == "" || newCommitID == "" {
		return nil, nil
	}

	opts := &issues_model.CreateCommentOptions{
		Type:        issues_model.CommentTypePullRequestPush,
		Doer:        pusher,
		Repo:        pr.BaseRepo,
		IsForcePush: isForcePush,
		Issue:       pr.Issue,
	}

	var data issues_model.PushActionContent
	if opts.IsForcePush {
		data.CommitIDs = []string{oldCommitID, newCommitID}
		data.IsForcePush = true
	} else {
		data.CommitIDs, err = getCommitIDsFromRepo(ctx, pr.BaseRepo, oldCommitID, newCommitID, pr.BaseBranch)
		if err != nil {
			return nil, err
		}
		// It maybe an empty pull request. Only non-empty pull request need to create push comment
		if len(data.CommitIDs) == 0 {
			return nil, nil
		}
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	opts.Content = string(dataJSON)
	comment, err = issues_model.CreateComment(ctx, opts)

	return comment, err
}
