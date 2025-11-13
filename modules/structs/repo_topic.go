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

// TopicResponse for returning topics
type TopicResponse struct {
	// The unique identifier of the topic
	ID int64 `json:"id"`
	// The name of the topic
	Name string `json:"topic_name"`
	// The number of repositories using this topic
	RepoCount int `json:"repo_count"`
	// The date and time when the topic was created
	Created time.Time `json:"created"`
	// The date and time when the topic was last updated
	Updated time.Time `json:"updated"`
}

// TopicName a list of repo topic names
type TopicName struct {
	// List of topic names
	TopicNames []string `json:"topics"`
}

// RepoTopicOptions a collection of repo topic names
type RepoTopicOptions struct {
	// list of topic names
	Topics []string `json:"topics"`
}
