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

	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddLanguageStats(x *xorm.Engine) error {
	// LanguageStat see models/repo_language_stats.go
	type LanguageStat struct {
		ID          int64 `xorm:"pk autoincr"`
		RepoID      int64 `xorm:"UNIQUE(s) INDEX NOT NULL"`
		CommitID    string
		IsPrimary   bool
		Language    string             `xorm:"VARCHAR(30) UNIQUE(s) INDEX NOT NULL"`
		Percentage  float32            `xorm:"NUMERIC(5,2) NOT NULL DEFAULT 0"`
		Color       string             `xorm:"-"`
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX CREATED"`
	}

	type RepoIndexerType int

	// RepoIndexerStatus see models/repo_stats_indexer.go
	type RepoIndexerStatus struct {
		ID          int64           `xorm:"pk autoincr"`
		RepoID      int64           `xorm:"INDEX(s)"`
		CommitSha   string          `xorm:"VARCHAR(40)"`
		IndexerType RepoIndexerType `xorm:"INDEX(s) NOT NULL DEFAULT 0"`
	}

	if err := x.Sync(new(LanguageStat)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}
	if err := x.Sync(new(RepoIndexerStatus)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}
	return nil
}
