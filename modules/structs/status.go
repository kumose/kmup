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

	"github.com/kumose/kmup/modules/commitstatus"
)

// CommitStatus holds a single status of a single Commit
type CommitStatus struct {
	// ID is the unique identifier for the commit status
	ID int64 `json:"id"`
	// State represents the status state (pending, success, error, failure)
	State commitstatus.CommitStatusState `json:"status"`
	// TargetURL is the URL to link to for more details
	TargetURL string `json:"target_url"`
	// Description provides a brief description of the status
	Description string `json:"description"`
	// URL is the API URL for this status
	URL string `json:"url"`
	// Context is the unique context identifier for the status
	Context string `json:"context"`
	// Creator is the user who created the status
	Creator *User `json:"creator"`
	// swagger:strfmt date-time
	Created time.Time `json:"created_at"`
	// swagger:strfmt date-time
	Updated time.Time `json:"updated_at"`
}

// CombinedStatus holds the combined state of several statuses for a single commit
type CombinedStatus struct {
	// State is the overall combined status state
	State commitstatus.CommitStatusState `json:"state"`
	// SHA is the commit SHA this status applies to
	SHA string `json:"sha"`
	// TotalCount is the total number of statuses
	TotalCount int `json:"total_count"`
	// Statuses contains all individual commit statuses
	Statuses []*CommitStatus `json:"statuses"`
	// Repository is the repository this status belongs to
	Repository *Repository `json:"repository"`
	// CommitURL is the API URL for the commit
	CommitURL string `json:"commit_url"`
	// URL is the API URL for this combined status
	URL string `json:"url"`
}

// CreateStatusOption holds the information needed to create a new CommitStatus for a Commit
type CreateStatusOption struct {
	// State represents the status state to set (pending, success, error, failure)
	State commitstatus.CommitStatusState `json:"state"`
	// TargetURL is the URL to link to for more details
	TargetURL string `json:"target_url"`
	// Description provides a brief description of the status
	Description string `json:"description"`
	// Context is the unique context identifier for the status
	Context string `json:"context"`
}
