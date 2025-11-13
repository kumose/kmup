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

package git

import (
	"context"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/commitstatus"
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/builder"
)

// CommitStatusSummary holds the latest commit Status of a single Commit
type CommitStatusSummary struct {
	ID        int64                          `xorm:"pk autoincr"`
	RepoID    int64                          `xorm:"INDEX UNIQUE(repo_id_sha)"`
	SHA       string                         `xorm:"VARCHAR(64) NOT NULL INDEX UNIQUE(repo_id_sha)"`
	State     commitstatus.CommitStatusState `xorm:"VARCHAR(7) NOT NULL"`
	TargetURL string                         `xorm:"TEXT"`
}

func init() {
	db.RegisterModel(new(CommitStatusSummary))
}

type RepoSHA struct {
	RepoID int64
	SHA    string
}

func GetLatestCommitStatusForRepoAndSHAs(ctx context.Context, repoSHAs []RepoSHA) ([]*CommitStatus, error) {
	cond := builder.NewCond()
	for _, rs := range repoSHAs {
		cond = cond.Or(builder.Eq{"repo_id": rs.RepoID, "sha": rs.SHA})
	}

	var summaries []CommitStatusSummary
	if err := db.GetEngine(ctx).Where(cond).Find(&summaries); err != nil {
		return nil, err
	}

	commitStatuses := make([]*CommitStatus, 0, len(repoSHAs))
	for _, summary := range summaries {
		commitStatuses = append(commitStatuses, &CommitStatus{
			RepoID:    summary.RepoID,
			SHA:       summary.SHA,
			State:     summary.State,
			TargetURL: summary.TargetURL,
		})
	}
	return commitStatuses, nil
}

func UpdateCommitStatusSummary(ctx context.Context, repoID int64, sha string) error {
	commitStatuses, err := GetLatestCommitStatus(ctx, repoID, sha, db.ListOptionsAll)
	if err != nil {
		return err
	}
	// it guarantees that commitStatuses is not empty because this function is always called after a commit status is created
	if len(commitStatuses) == 0 {
		setting.PanicInDevOrTesting("no commit statuses found for repo %d and sha %s", repoID, sha)
	}
	state := CalcCommitStatus(commitStatuses) // non-empty commitStatuses is guaranteed
	// mysql will return 0 when update a record which state hasn't been changed which behaviour is different from other database,
	// so we need to use insert in on duplicate
	if setting.Database.Type.IsMySQL() {
		_, err := db.GetEngine(ctx).Exec("INSERT INTO commit_status_summary (repo_id,sha,state,target_url) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE state=?",
			repoID, sha, state.State, state.TargetURL, state.State)
		return err
	}

	if cnt, err := db.GetEngine(ctx).Where("repo_id=? AND sha=?", repoID, sha).
		Cols("state, target_url").
		Update(&CommitStatusSummary{
			State:     state.State,
			TargetURL: state.TargetURL,
		}); err != nil {
		return err
	} else if cnt == 0 {
		_, err = db.GetEngine(ctx).Insert(&CommitStatusSummary{
			RepoID:    repoID,
			SHA:       sha,
			State:     state.State,
			TargetURL: state.TargetURL,
		})
		return err
	}
	return nil
}
