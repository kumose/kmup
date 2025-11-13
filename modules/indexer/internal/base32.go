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

package internal

import (
	"fmt"
	"strconv"
)

func Base36(i int64) string {
	return strconv.FormatInt(i, 36)
}

func ParseBase36(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 36, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid base36 integer %q: %w", s, err)
	}
	return i, nil
}
