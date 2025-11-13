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

package context

import (
	"strings"
	"time"
)

// GetQueryBeforeSince return parsed time (unix format) from URL query's before and since
func GetQueryBeforeSince(ctx *Base) (before, since int64, err error) {
	before, err = parseFormTime(ctx, "before")
	if err != nil {
		return 0, 0, err
	}

	since, err = parseFormTime(ctx, "since")
	if err != nil {
		return 0, 0, err
	}
	return before, since, nil
}

// parseTime parse time and return unix timestamp
func parseFormTime(ctx *Base, name string) (int64, error) {
	value := strings.TrimSpace(ctx.FormString(name))
	if len(value) != 0 {
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return 0, err
		}
		if !t.IsZero() {
			return t.Unix(), nil
		}
	}
	return 0, nil
}
