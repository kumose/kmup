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

// GitBlobResponse represents a git blob
type GitBlobResponse struct {
	// The content of the git blob (may be base64 encoded)
	Content *string `json:"content"`
	// The encoding used for the content (e.g., "base64")
	Encoding *string `json:"encoding"`
	// The URL to access this git blob
	URL string `json:"url"`
	// The SHA hash of the git blob
	SHA string `json:"sha"`
	// The size of the git blob in bytes
	Size int64 `json:"size"`

	// The LFS object ID if this blob is stored in LFS
	LfsOid *string `json:"lfs_oid,omitempty"`
	// The size of the LFS object if this blob is stored in LFS
	LfsSize *int64 `json:"lfs_size,omitempty"`
}
