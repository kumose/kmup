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

// WatchInfo represents an API watch status of one repository
type WatchInfo struct {
	// Whether the repository is being watched for notifications
	Subscribed bool `json:"subscribed"`
	// Whether notifications for the repository are ignored
	Ignored bool `json:"ignored"`
	// The reason for the current watch status
	Reason any `json:"reason"`
	// The timestamp when the watch status was created
	CreatedAt time.Time `json:"created_at"`
	// The URL for managing the watch status
	URL string `json:"url"`
	// The URL of the repository being watched
	RepositoryURL string `json:"repository_url"`
}
