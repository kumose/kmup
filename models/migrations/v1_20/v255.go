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
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddArchivedUnixToRepository(x *xorm.Engine) error {
	type Repository struct {
		ArchivedUnix timeutil.TimeStamp `xorm:"DEFAULT 0"`
	}

	if err := x.Sync(new(Repository)); err != nil {
		return err
	}

	_, err := x.Exec("UPDATE repository SET archived_unix = updated_unix WHERE is_archived = ? AND archived_unix = 0", true)
	return err
}
