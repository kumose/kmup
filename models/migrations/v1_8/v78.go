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

package v1_8

import (
	"github.com/kumose/kmup/models/migrations/base"

	"xorm.io/xorm"
)

func RenameRepoIsBareToIsEmpty(x *xorm.Engine) error {
	type Repository struct {
		ID      int64 `xorm:"pk autoincr"`
		IsBare  bool
		IsEmpty bool `xorm:"INDEX"`
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := sess.Sync(new(Repository)); err != nil {
		return err
	}
	if _, err := sess.Exec("UPDATE repository SET is_empty = is_bare;"); err != nil {
		return err
	}
	if err := sess.Commit(); err != nil {
		return err
	}

	if err := sess.Begin(); err != nil {
		return err
	}
	if err := base.DropTableColumns(sess, "repository", "is_bare"); err != nil {
		return err
	}

	return sess.Commit()
}
