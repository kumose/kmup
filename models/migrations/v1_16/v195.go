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

package v1_16

import (
	"fmt"

	"xorm.io/xorm"
)

func AddTableCommitStatusIndex(x *xorm.Engine) error {
	// CommitStatusIndex represents a table for commit status index
	type CommitStatusIndex struct {
		ID       int64
		RepoID   int64  `xorm:"unique(repo_sha)"`
		SHA      string `xorm:"unique(repo_sha)"`
		MaxIndex int64  `xorm:"index"`
	}

	if err := x.Sync(new(CommitStatusIndex)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	// Remove data we're goint to rebuild
	if _, err := sess.Table("commit_status_index").Where("1=1").Delete(&CommitStatusIndex{}); err != nil {
		return err
	}

	// Create current data for all repositories with issues and PRs
	if _, err := sess.Exec("INSERT INTO commit_status_index (repo_id, sha, max_index) " +
		"SELECT max_data.repo_id, max_data.sha, max_data.max_index " +
		"FROM ( SELECT commit_status.repo_id AS repo_id, commit_status.sha AS sha, max(commit_status.`index`) AS max_index " +
		"FROM commit_status GROUP BY commit_status.repo_id, commit_status.sha) AS max_data"); err != nil {
		return err
	}

	return sess.Commit()
}
