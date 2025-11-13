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

func FixIncorrectOwnerTeamUnitAccessMode(x *xorm.Engine) error {
	type UnitType int
	type AccessMode int

	type TeamUnit struct {
		ID         int64    `xorm:"pk autoincr"`
		OrgID      int64    `xorm:"INDEX"`
		TeamID     int64    `xorm:"UNIQUE(s)"`
		Type       UnitType `xorm:"UNIQUE(s)"`
		AccessMode AccessMode
	}

	const (
		// AccessModeOwner owner access
		AccessModeOwner = 4
	)

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	count, err := sess.Table("team_unit").
		Where("team_id IN (SELECT id FROM team WHERE authorize = ?)", AccessModeOwner).
		Update(&TeamUnit{
			AccessMode: AccessModeOwner,
		})
	if err != nil {
		return err
	}
	log.Debug("Updated %d owner team unit access mode to belong to owner instead of none", count)

	return sess.Commit()
}
