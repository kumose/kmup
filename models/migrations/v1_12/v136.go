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

package v1_12

import (
	"fmt"
	"math"
	"time"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm"
)

func AddCommitDivergenceToPulls(x *xorm.Engine) error {
	type Repository struct {
		ID        int64 `xorm:"pk autoincr"`
		OwnerID   int64 `xorm:"UNIQUE(s) index"`
		OwnerName string
		LowerName string `xorm:"UNIQUE(s) INDEX NOT NULL"`
		Name      string `xorm:"INDEX NOT NULL"`
	}

	type PullRequest struct {
		ID      int64 `xorm:"pk autoincr"`
		IssueID int64 `xorm:"INDEX"`
		Index   int64

		CommitsAhead  int
		CommitsBehind int

		BaseRepoID int64 `xorm:"INDEX"`
		BaseBranch string

		HasMerged      bool   `xorm:"INDEX"`
		MergedCommitID string `xorm:"VARCHAR(40)"`
	}

	if err := x.Sync(new(PullRequest)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	last := 0
	migrated := 0

	batchSize := setting.Database.IterateBufferSize
	sess := x.NewSession()
	defer sess.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	count, err := sess.Where("has_merged = ?", false).Count(new(PullRequest))
	if err != nil {
		return err
	}
	log.Info("%d Unmerged Pull Request(s) to migrate ...", count)

	for {
		if err := sess.Begin(); err != nil {
			return err
		}
		results := make([]*PullRequest, 0, batchSize)
		err := sess.Where("has_merged = ?", false).OrderBy("id").Limit(batchSize, last).Find(&results)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			break
		}
		last += batchSize

		for _, pr := range results {
			baseRepo := &Repository{ID: pr.BaseRepoID}
			has, err := x.Table("repository").Get(baseRepo)
			if err != nil {
				return fmt.Errorf("Unable to get base repo %d %w", pr.BaseRepoID, err)
			}
			if !has {
				log.Error("Missing base repo with id %d for PR ID %d", pr.BaseRepoID, pr.ID)
				continue
			}
			repoStore := repo_model.StorageRepo(repo_model.RelativePath(baseRepo.OwnerName, baseRepo.Name))
			gitRefName := fmt.Sprintf("refs/pull/%d/head", pr.Index)
			divergence, err := gitrepo.GetDivergingCommits(graceful.GetManager().HammerContext(), repoStore, pr.BaseBranch, gitRefName)
			if err != nil {
				log.Warn("Could not recalculate Divergence for pull: %d", pr.ID)
				pr.CommitsAhead = 0
				pr.CommitsBehind = 0
			}
			pr.CommitsAhead = divergence.Ahead
			pr.CommitsBehind = divergence.Behind

			if _, err = sess.ID(pr.ID).Cols("commits_ahead", "commits_behind").Update(pr); err != nil {
				return fmt.Errorf("Update Cols: %w", err)
			}
			migrated++
		}

		if err := sess.Commit(); err != nil {
			return err
		}
		select {
		case <-ticker.C:
			log.Info(
				"%d/%d (%2.0f%%) Pull Request(s) migrated in %d batches. %d PRs Remaining ...",
				migrated,
				count,
				float64(migrated)/float64(count)*100,
				int(math.Ceil(float64(migrated)/float64(batchSize))),
				count-int64(migrated))
		default:
		}
	}
	log.Info("Completed migrating %d Pull Request(s) in: %d batches", count, int(math.Ceil(float64(migrated)/float64(batchSize))))
	return nil
}
