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

package timeutil

import (
	"time"

	"github.com/kumose/kmup/modules/setting"
)

// TimeStampNano is for nano time in database, do not use it unless there is a real requirement.
type TimeStampNano int64

// TimeStampNanoNow returns now nano int64
func TimeStampNanoNow() TimeStampNano {
	return TimeStampNano(time.Now().UnixNano())
}

// AsTime convert timestamp as time.Time in Local locale
func (tsn TimeStampNano) AsTime() (tm time.Time) {
	return tsn.AsTimeInLocation(setting.DefaultUILocation)
}

// AsTimeInLocation convert timestamp as time.Time in Local locale
func (tsn TimeStampNano) AsTimeInLocation(loc *time.Location) time.Time {
	return time.Unix(0, int64(tsn)).In(loc)
}
