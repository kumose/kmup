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

// GitEntry represents a git tree
type GitEntry struct {
	// Path is the file or directory path
	Path string `json:"path"`
	// Mode is the file mode (permissions)
	Mode string `json:"mode"`
	// Type indicates if this is a file, directory, or symlink
	Type string `json:"type"`
	// Size is the file size in bytes
	Size int64 `json:"size"`
	// SHA is the Git object SHA
	SHA string `json:"sha"`
	// URL is the API URL for this tree entry
	URL string `json:"url"`
}

// GitTreeResponse returns a git tree
type GitTreeResponse struct {
	// SHA is the tree object SHA
	SHA string `json:"sha"`
	// URL is the API URL for this tree
	URL string `json:"url"`
	// Entries contains the tree entries (files and directories)
	Entries []GitEntry `json:"tree"`
	// Truncated indicates if the response was truncated due to size
	Truncated bool `json:"truncated"`
	// Page is the current page number for pagination
	Page int `json:"page"`
	// TotalCount is the total number of entries in the tree
	TotalCount int `json:"total_count"`
}
