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

package migration

import "github.com/kumose/kmup/modules/structs"

// MigrateOptions defines the way a repository gets migrated
// this is for internal usage by migrations module and func who interact with it
type MigrateOptions struct {
	// required: true
	CloneAddr             string `json:"clone_addr" binding:"Required"`
	CloneAddrEncrypted    string `json:"clone_addr_encrypted,omitempty"`
	AuthUsername          string `json:"auth_username"`
	AuthPassword          string `json:"-"`
	AuthPasswordEncrypted string `json:"auth_password_encrypted,omitempty"`
	AuthToken             string `json:"-"`
	AuthTokenEncrypted    string `json:"auth_token_encrypted,omitempty"`
	// required: true
	UID int `json:"uid" binding:"Required"`
	// required: true
	RepoName        string `json:"repo_name" binding:"Required"`
	Mirror          bool   `json:"mirror"`
	LFS             bool   `json:"lfs"`
	LFSEndpoint     string `json:"lfs_endpoint"`
	Private         bool   `json:"private"`
	Description     string `json:"description"`
	OriginalURL     string
	GitServiceType  structs.GitServiceType
	Wiki            bool
	Issues          bool
	Milestones      bool
	Labels          bool
	Releases        bool
	Comments        bool
	PullRequests    bool
	ReleaseAssets   bool
	MigrateToRepoID int64
	MirrorInterval  string `json:"mirror_interval"`

	AWSAccessKeyID     string
	AWSSecretAccessKey string
}
