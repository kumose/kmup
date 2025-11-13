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
	"fmt"
	"strconv"
	"time"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/graceful"
	logger "github.com/kumose/kmup/modules/log"
)

// CompareInfo represents needed information for comparing references.
type CompareInfo struct {
	MergeBase    string
	BaseCommitID string
	HeadCommitID string
	Commits      []*git.Commit
	NumFiles     int
}

// GetCompareInfo generates and returns compare information between base and head branches of repositories.
func GetCompareInfo(ctx context.Context, baseRepo, headRepo *repo_model.Repository, headGitRepo *git.Repository, baseBranch, headBranch string, directComparison, fileOnly bool) (_ *CompareInfo, err error) {
	var (
		remoteBranch string
		tmpRemote    string
	)

	// We don't need a temporary remote for same repository.
	if headGitRepo.Path != baseRepo.RepoPath() {
		// Add a temporary remote
		tmpRemote = strconv.FormatInt(time.Now().UnixNano(), 10)
		if err = gitrepo.GitRemoteAdd(ctx, headRepo, tmpRemote, baseRepo.RepoPath()); err != nil {
			return nil, fmt.Errorf("GitRemoteAdd: %w", err)
		}
		defer func() {
			if err := gitrepo.GitRemoteRemove(graceful.GetManager().ShutdownContext(), headRepo, tmpRemote); err != nil {
				logger.Error("GetPullRequestInfo: GitRemoteRemove: %v", err)
			}
		}()
	}

	compareInfo := new(CompareInfo)

	compareInfo.HeadCommitID, err = git.GetFullCommitID(ctx, headGitRepo.Path, headBranch)
	if err != nil {
		compareInfo.HeadCommitID = headBranch
	}

	compareInfo.MergeBase, remoteBranch, err = headGitRepo.GetMergeBase(tmpRemote, baseBranch, headBranch)
	if err == nil {
		compareInfo.BaseCommitID, err = git.GetFullCommitID(ctx, headGitRepo.Path, remoteBranch)
		if err != nil {
			compareInfo.BaseCommitID = remoteBranch
		}
		separator := "..."
		baseCommitID := compareInfo.MergeBase
		if directComparison {
			separator = ".."
			baseCommitID = compareInfo.BaseCommitID
		}

		// We have a common base - therefore we know that ... should work
		if !fileOnly {
			compareInfo.Commits, err = headGitRepo.ShowPrettyFormatLogToList(ctx, baseCommitID+separator+headBranch)
			if err != nil {
				return nil, fmt.Errorf("ShowPrettyFormatLogToList: %w", err)
			}
		} else {
			compareInfo.Commits = []*git.Commit{}
		}
	} else {
		compareInfo.Commits = []*git.Commit{}
		compareInfo.MergeBase, err = git.GetFullCommitID(ctx, headGitRepo.Path, remoteBranch)
		if err != nil {
			compareInfo.MergeBase = remoteBranch
		}
		compareInfo.BaseCommitID = compareInfo.MergeBase
	}

	// Count number of changed files.
	// This probably should be removed as we need to use shortstat elsewhere
	// Now there is git diff --shortstat but this appears to be slower than simply iterating with --nameonly
	compareInfo.NumFiles, err = headGitRepo.GetDiffNumChangedFiles(remoteBranch, headBranch, directComparison)
	if err != nil {
		return nil, err
	}
	return compareInfo, nil
}
