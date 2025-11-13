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

package v1_22

import "xorm.io/xorm"

func AddCommitStatusSummary(x *xorm.Engine) error {
	type CommitStatusSummary struct {
		ID     int64  `xorm:"pk autoincr"`
		RepoID int64  `xorm:"INDEX UNIQUE(repo_id_sha)"`
		SHA    string `xorm:"VARCHAR(64) NOT NULL INDEX UNIQUE(repo_id_sha)"`
		State  string `xorm:"VARCHAR(7) NOT NULL"`
	}
	// there is no migrations because if there is no data on this table, it will fall back to get data
	// from commit status
	return x.Sync2(new(CommitStatusSummary))
}
