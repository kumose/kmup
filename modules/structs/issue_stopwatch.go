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

// StopWatch represent a running stopwatch
type StopWatch struct {
	// swagger:strfmt date-time
	// Created is the time when the stopwatch was started
	Created time.Time `json:"created"`
	// Seconds is the total elapsed time in seconds
	Seconds int64 `json:"seconds"`
	// Duration is a human-readable duration string
	Duration string `json:"duration"`
	// IssueIndex is the index number of the associated issue
	IssueIndex int64 `json:"issue_index"`
	// IssueTitle is the title of the associated issue
	IssueTitle string `json:"issue_title"`
	// RepoOwnerName is the name of the repository owner
	RepoOwnerName string `json:"repo_owner_name"`
	// RepoName is the name of the repository
	RepoName string `json:"repo_name"`
}

// StopWatches represent a list of stopwatches
type StopWatches []StopWatch
