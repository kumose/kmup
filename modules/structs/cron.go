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

package structs

import "time"

// Cron represents a Cron task
type Cron struct {
	// The name of the cron task
	Name string `json:"name"`
	// The cron schedule expression (e.g., "0 0 * * *")
	Schedule string `json:"schedule"`
	// The next scheduled execution time
	Next time.Time `json:"next"`
	// The previous execution time
	Prev time.Time `json:"prev"`
	// The total number of times this cron task has been executed
	ExecTimes int64 `json:"exec_times"`
}
