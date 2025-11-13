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

package v1_18

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddTeamInviteTable(x *xorm.Engine) error {
	type TeamInvite struct {
		ID          int64              `xorm:"pk autoincr"`
		Token       string             `xorm:"UNIQUE(token) INDEX NOT NULL DEFAULT ''"`
		InviterID   int64              `xorm:"NOT NULL DEFAULT 0"`
		OrgID       int64              `xorm:"INDEX NOT NULL DEFAULT 0"`
		TeamID      int64              `xorm:"UNIQUE(team_mail) INDEX NOT NULL DEFAULT 0"`
		Email       string             `xorm:"UNIQUE(team_mail) NOT NULL DEFAULT ''"`
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	return x.Sync(new(TeamInvite))
}
