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

package convert

import (
	"strings"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/structs"
)

// ToCorrectPageSize makes sure page size is in allowed range.
func ToCorrectPageSize(size int) int {
	if size <= 0 {
		size = setting.API.DefaultPagingNum
	} else if size > setting.API.MaxResponseItems {
		size = setting.API.MaxResponseItems
	}
	return size
}

// ToGitServiceType return GitServiceType based on string
func ToGitServiceType(value string) structs.GitServiceType {
	switch strings.ToLower(value) {
	case "github":
		return structs.GithubService
	case "kmup":
		return structs.KmupService
	case "gitlab":
		return structs.GitlabService
	case "gogs":
		return structs.GogsService
	case "onedev":
		return structs.OneDevService
	case "gitbucket":
		return structs.GitBucketService
	case "codebase":
		return structs.CodebaseService
	case "codecommit":
		return structs.CodeCommitService
	default:
		return structs.PlainGitService
	}
}
