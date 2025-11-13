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

package organization

import (
	"context"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/perm"
	"github.com/kumose/kmup/models/unit"
)

// TeamUnit describes all units of a repository
type TeamUnit struct {
	ID         int64     `xorm:"pk autoincr"`
	OrgID      int64     `xorm:"INDEX"`
	TeamID     int64     `xorm:"UNIQUE(s)"`
	Type       unit.Type `xorm:"UNIQUE(s)"`
	AccessMode perm.AccessMode
}

// Unit returns Unit
func (t *TeamUnit) Unit() unit.Unit {
	return unit.Units[t.Type]
}

func getUnitsByTeamID(ctx context.Context, teamID int64) (units []*TeamUnit, err error) {
	return units, db.GetEngine(ctx).Where("team_id = ?", teamID).Find(&units)
}

// UpdateTeamUnits updates a teams's units
func UpdateTeamUnits(ctx context.Context, team *Team, units []TeamUnit) (err error) {
	return db.WithTx(ctx, func(ctx context.Context) error {
		if _, err = db.GetEngine(ctx).Where("team_id = ?", team.ID).Delete(new(TeamUnit)); err != nil {
			return err
		}

		if len(units) > 0 {
			if err = db.Insert(ctx, units); err != nil {
				return err
			}
		}
		return nil
	})
}
