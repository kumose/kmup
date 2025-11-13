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

import (
	"xorm.io/xorm"
)

type BadgeUnique struct {
	ID   int64  `xorm:"pk autoincr"`
	Slug string `xorm:"UNIQUE"`
}

func (BadgeUnique) TableName() string {
	return "badge"
}

func UseSlugInsteadOfIDForBadges(x *xorm.Engine) error {
	type Badge struct {
		Slug string
	}

	err := x.Sync(new(Badge))
	if err != nil {
		return err
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	_, err = sess.Exec("UPDATE `badge` SET `slug` = `id` Where `slug` IS NULL")
	if err != nil {
		return err
	}

	err = sess.Sync(new(BadgeUnique))
	if err != nil {
		return err
	}

	return sess.Commit()
}
