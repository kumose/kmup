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

package v1_26

import (
	"xorm.io/xorm"
)

func AddActionsConcurrency(x *xorm.Engine) error {
	type ActionRun struct {
		RepoID            int64 `xorm:"index(repo_concurrency)"`
		RawConcurrency    string
		ConcurrencyGroup  string `xorm:"index(repo_concurrency) NOT NULL DEFAULT ''"`
		ConcurrencyCancel bool   `xorm:"NOT NULL DEFAULT FALSE"`
	}

	if _, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreDropIndices: true,
	}, new(ActionRun)); err != nil {
		return err
	}

	if err := x.Sync(new(ActionRun)); err != nil {
		return err
	}

	type ActionRunJob struct {
		RepoID                 int64 `xorm:"index(repo_concurrency)"`
		RawConcurrency         string
		IsConcurrencyEvaluated bool
		ConcurrencyGroup       string `xorm:"index(repo_concurrency) NOT NULL DEFAULT ''"`
		ConcurrencyCancel      bool   `xorm:"NOT NULL DEFAULT FALSE"`
	}

	if _, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreDropIndices: true,
	}, new(ActionRunJob)); err != nil {
		return err
	}

	return nil
}
