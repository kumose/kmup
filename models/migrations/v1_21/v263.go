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

package v1_21

import (
	"fmt"

	"xorm.io/xorm"
)

// AddGitSizeAndLFSSizeToRepositoryTable: add GitSize and LFSSize columns to Repository
func AddGitSizeAndLFSSizeToRepositoryTable(x *xorm.Engine) error {
	type Repository struct {
		GitSize int64 `xorm:"NOT NULL DEFAULT 0"`
		LFSSize int64 `xorm:"NOT NULL DEFAULT 0"`
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if err := sess.Sync(new(Repository)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	_, err := sess.Exec(`UPDATE repository SET lfs_size=(SELECT SUM(size) FROM lfs_meta_object WHERE lfs_meta_object.repository_id=repository.ID) WHERE EXISTS (SELECT 1 FROM lfs_meta_object WHERE lfs_meta_object.repository_id=repository.ID)`)
	if err != nil {
		return err
	}

	_, err = sess.Exec(`UPDATE repository SET size = 0 WHERE size IS NULL`)
	if err != nil {
		return err
	}

	_, err = sess.Exec(`UPDATE repository SET git_size = size - lfs_size WHERE size > lfs_size`)
	if err != nil {
		return err
	}

	return sess.Commit()
}
