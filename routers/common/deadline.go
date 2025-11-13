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

package common

import (
	"time"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/timeutil"
)

func ParseDeadlineDateToEndOfDay(date string) (timeutil.TimeStamp, error) {
	if date == "" {
		return 0, nil
	}
	deadline, err := time.ParseInLocation("2006-01-02", date, setting.DefaultUILocation)
	if err != nil {
		return 0, err
	}
	deadline = time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 23, 59, 59, 0, deadline.Location())
	return timeutil.TimeStamp(deadline.Unix()), nil
}

func ParseAPIDeadlineToEndOfDay(t *time.Time) (timeutil.TimeStamp, error) {
	if t == nil || t.IsZero() || t.Unix() == 0 {
		return 0, nil
	}
	deadline := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, setting.DefaultUILocation)
	return timeutil.TimeStamp(deadline.Unix()), nil
}
