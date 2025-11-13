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

import (
	"time"
)

// AddTimeOption options for adding time to an issue
type AddTimeOption struct {
	// time in seconds
	// required: true
	Time int64 `json:"time" binding:"Required"`
	// swagger:strfmt date-time
	Created time.Time `json:"created"`
	// username of the user who spent the time working on the issue (optional)
	User string `json:"user_name"`
}

// TrackedTime worked time for an issue / pr
type TrackedTime struct {
	// ID is the unique identifier for the tracked time entry
	ID int64 `json:"id"`
	// swagger:strfmt date-time
	Created time.Time `json:"created"`
	// Time in seconds
	Time int64 `json:"time"`
	// deprecated (only for backwards compatibility)
	UserID int64 `json:"user_id"`
	// username of the user
	UserName string `json:"user_name"`
	// deprecated (only for backwards compatibility)
	IssueID int64 `json:"issue_id"`
	// Issue contains the associated issue information
	Issue *Issue `json:"issue"`
}

// TrackedTimeList represents a list of tracked times
type TrackedTimeList []*TrackedTime
