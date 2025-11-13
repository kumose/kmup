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

package v1_20

import (
	"github.com/kumose/kmup/modules/log"

	"xorm.io/xorm"
)

// FixIncorrectProjectType: set individual project's type from 3(TypeOrganization) to 1(TypeIndividual)
func FixIncorrectProjectType(x *xorm.Engine) error {
	type User struct {
		ID   int64 `xorm:"pk autoincr"`
		Type int
	}

	const (
		UserTypeIndividual int = 0

		TypeIndividual   uint8 = 1
		TypeOrganization uint8 = 3
	)

	type Project struct {
		OwnerID int64 `xorm:"INDEX"`
		Type    uint8
		Owner   *User `xorm:"extends"`
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	count, err := sess.Table("project").
		Where("type = ? AND owner_id IN (SELECT id FROM `user` WHERE type = ?)", TypeOrganization, UserTypeIndividual).
		Update(&Project{
			Type: TypeIndividual,
		})
	if err != nil {
		return err
	}
	log.Debug("Updated %d projects to belong to a user instead of an organization", count)

	return sess.Commit()
}
