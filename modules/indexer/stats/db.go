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

package stats

import (
	"fmt"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/languagestats"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/setting"
)

// DBIndexer implements Indexer interface to use database's like search
type DBIndexer struct{}

// Index repository status function
func (db *DBIndexer) Index(id int64) error {
	ctx, _, finished := process.GetManager().AddContext(graceful.GetManager().ShutdownContext(), fmt.Sprintf("Stats.DB Index Repo[%d]", id))
	defer finished()

	repo, err := repo_model.GetRepositoryByID(ctx, id)
	if err != nil {
		return err
	}
	if repo.IsEmpty {
		return nil
	}

	status, err := repo_model.GetIndexerStatus(ctx, repo, repo_model.RepoIndexerTypeStats)
	if err != nil {
		return err
	}

	gitRepo, err := gitrepo.OpenRepository(ctx, repo)
	if err != nil {
		if err.Error() == "no such file or directory" {
			return nil
		}
		return err
	}
	defer gitRepo.Close()

	// Get latest commit for default branch
	commitID, err := gitRepo.GetBranchCommitID(repo.DefaultBranch)
	if err != nil {
		if git.IsErrBranchNotExist(err) || git.IsErrNotExist(err) || setting.IsInTesting {
			log.Debug("Unable to get commit ID for default branch %s in %s ... skipping this repository", repo.DefaultBranch, repo.FullName())
			return nil
		}
		log.Error("Unable to get commit ID for default branch %s in %s. Error: %v", repo.DefaultBranch, repo.FullName(), err)
		return err
	}

	// Do not recalculate stats if already calculated for this commit
	if status.CommitSha == commitID {
		return nil
	}

	// Calculate and save language statistics to database
	stats, err := languagestats.GetLanguageStats(gitRepo, commitID)
	if err != nil {
		if !setting.IsInTesting {
			log.Error("Unable to get language stats for ID %s for default branch %s in %s. Error: %v", commitID, repo.DefaultBranch, repo.FullName(), err)
		}
		return err
	}
	err = repo_model.UpdateLanguageStats(ctx, repo, commitID, stats)
	if err != nil {
		log.Error("Unable to update language stats for ID %s for default branch %s in %s. Error: %v", commitID, repo.DefaultBranch, repo.FullName(), err)
		return err
	}

	log.Debug("DBIndexer completed language stats for ID %s for default branch %s in %s. stats count: %d", commitID, repo.DefaultBranch, repo.FullName(), len(stats))
	return nil
}

// Close dummy function
func (db *DBIndexer) Close() {
}
