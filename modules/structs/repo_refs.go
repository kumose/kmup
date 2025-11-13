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

// Reference represents a Git reference.
type Reference struct {
	// The name of the Git reference (e.g., refs/heads/master)
	Ref string `json:"ref"`
	// The URL to access this Git reference
	URL string `json:"url"`
	// The Git object that this reference points to
	Object *GitObject `json:"object"`
}

// GitObject represents a Git object.
type GitObject struct {
	// The type of the Git object (e.g., commit, tag, tree, blob)
	Type string `json:"type"`
	// The SHA hash of the Git object
	SHA string `json:"sha"`
	// The URL to access this Git object
	URL string `json:"url"`
}
