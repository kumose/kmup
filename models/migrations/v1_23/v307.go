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

package v1_23

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func FixMilestoneNoDueDate(x *xorm.Engine) error {
	type Milestone struct {
		DeadlineUnix timeutil.TimeStamp
	}
	// Wednesday, December 1, 9999 12:00:00 AM GMT+00:00
	_, err := x.Table("milestone").Where("deadline_unix > 253399622400").
		Cols("deadline_unix").
		Update(&Milestone{DeadlineUnix: 0})
	return err
}
