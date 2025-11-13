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
	"github.com/kumose/kmup/modules/log"

	"xorm.io/xorm"
)

func RecalculateStars(x *xorm.Engine) (err error) {

	// recalculate Stars number for all users to fully fix it.

	type User struct {
		ID int64 `xorm:"pk autoincr"`
	}

	const batchSize = 100
	sess := x.NewSession()
	defer sess.Close()

	for start := 0; ; start += batchSize {
		users := make([]User, 0, batchSize)
		if err := sess.Limit(batchSize, start).Where("type = ?", 0).Cols("id").Find(&users); err != nil {
			return err
		}
		if len(users) == 0 {
			break
		}

		if err := sess.Begin(); err != nil {
			return err
		}

		for _, user := range users {
			if _, err := sess.Exec("UPDATE `user` SET num_stars=(SELECT COUNT(*) FROM `star` WHERE uid=?) WHERE id=?", user.ID, user.ID); err != nil {
				return err
			}
		}

		if err := sess.Commit(); err != nil {
			return err
		}
	}

	log.Debug("recalculate Stars number for all user finished")

	return err
}
