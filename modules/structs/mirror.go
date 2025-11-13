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

// CreatePushMirrorOption represents need information to create a push mirror of a repository.
type CreatePushMirrorOption struct {
	// The remote repository URL to push to
	RemoteAddress string `json:"remote_address"`
	// The username for authentication with the remote repository
	RemoteUsername string `json:"remote_username"`
	// The password for authentication with the remote repository
	RemotePassword string `json:"remote_password"`
	// The sync interval for automatic updates
	Interval string `json:"interval"`
	// Whether to sync on every commit
	SyncOnCommit bool `json:"sync_on_commit"`
}

// PushMirror represents information of a push mirror
// swagger:model
type PushMirror struct {
	// The name of the source repository
	RepoName string `json:"repo_name"`
	// The name of the remote in the git configuration
	RemoteName string `json:"remote_name"`
	// The remote repository URL being mirrored to
	RemoteAddress string `json:"remote_address"`
	// swagger:strfmt date-time
	CreatedUnix time.Time `json:"created"`
	// swagger:strfmt date-time
	LastUpdateUnix *time.Time `json:"last_update"`
	// The last error message encountered during sync
	LastError string `json:"last_error"`
	// The sync interval for automatic updates
	Interval string `json:"interval"`
	// Whether to sync on every commit
	SyncOnCommit bool `json:"sync_on_commit"`
}
