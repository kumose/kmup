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
	"math"

	"github.com/kumose/kmup/models/db"
)

// ParsePaginator parses a db.Paginator into a skip and limit
func ParsePaginator(paginator *db.ListOptions, maxNums ...int) (int, int) {
	// Use a very large number to indicate no limit
	unlimited := math.MaxInt32
	if len(maxNums) > 0 {
		// Some indexer engines have a limit on the page size, respect that
		unlimited = maxNums[0]
	}

	if paginator == nil || paginator.IsListAll() {
		// It shouldn't happen. In actual usage scenarios, there should not be requests to search all.
		// But if it does happen, respect it and return "unlimited".
		// And it's also useful for testing.
		return 0, unlimited
	}

	if paginator.PageSize == 0 {
		// Do not return any results when searching, it's used to get the total count only.
		return 0, 0
	}

	return paginator.GetSkipTake()
}
