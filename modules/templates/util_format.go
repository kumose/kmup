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

package templates

import (
	"fmt"
	"strconv"

	"github.com/kumose/kmup/modules/util"
)

func timeEstimateString(timeSec any) string {
	v, _ := util.ToInt64(timeSec)
	if v == 0 {
		return ""
	}
	return util.TimeEstimateString(v)
}

func countFmt(data any) string {
	// legacy code, not ideal, still used in some places
	num, err := util.ToInt64(data)
	if err != nil {
		return ""
	}
	if num < 1000 {
		return strconv.FormatInt(num, 10)
	} else if num < 1_000_000 {
		num2 := float32(num) / 1000.0
		return fmt.Sprintf("%.1fk", num2)
	} else if num < 1_000_000_000 {
		num2 := float32(num) / 1_000_000.0
		return fmt.Sprintf("%.1fM", num2)
	}
	num2 := float32(num) / 1_000_000_000.0
	return fmt.Sprintf("%.1fG", num2)
}
