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

package v1_13

import (
	"fmt"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm"
)

func FixLanguageStatsToSaveSize(x *xorm.Engine) error {
	// LanguageStat see models/repo_language_stats.go
	type LanguageStat struct {
		Size int64 `xorm:"NOT NULL DEFAULT 0"`
	}

	// RepoIndexerType specifies the repository indexer type
	type RepoIndexerType int

	const RepoIndexerTypeStats RepoIndexerType = 1

	// RepoIndexerStatus see models/repo_indexer.go
	type RepoIndexerStatus struct {
		IndexerType RepoIndexerType `xorm:"INDEX(s) NOT NULL DEFAULT 0"`
	}

	if err := x.Sync(new(LanguageStat)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	x.Delete(&RepoIndexerStatus{IndexerType: RepoIndexerTypeStats})

	// Delete language stat statuses
	truncExpr := "TRUNCATE TABLE"
	if setting.Database.Type.IsSQLite3() {
		truncExpr = "DELETE FROM"
	}

	// Delete language stats
	if _, err := x.Exec(truncExpr + " language_stat"); err != nil {
		return err
	}

	sess := x.NewSession()
	defer sess.Close()
	return base.DropTableColumns(sess, "language_stat", "percentage")
}
