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

// WikiCommit page commit/revision
type WikiCommit struct {
	// The commit SHA hash
	ID string `json:"sha"`
	// The author of the commit
	Author *CommitUser `json:"author"`
	// The committer of the commit
	Committer *CommitUser `json:"commiter"`
	// The commit message
	Message string `json:"message"`
}

// WikiPage a wiki page
type WikiPage struct {
	*WikiPageMetaData
	// Page content, base64 encoded
	ContentBase64 string `json:"content_base64"`
	// The number of commits that modified this page
	CommitCount int64 `json:"commit_count"`
	// The sidebar content for the wiki page
	Sidebar string `json:"sidebar"`
	// The footer content for the wiki page
	Footer string `json:"footer"`
}

// WikiPageMetaData wiki page meta information
type WikiPageMetaData struct {
	// The title of the wiki page
	Title string `json:"title"`
	// The HTML URL to view the wiki page
	HTMLURL string `json:"html_url"`
	// The sub URL path for the wiki page
	SubURL string `json:"sub_url"`
	// The last commit that modified this wiki page
	LastCommit *WikiCommit `json:"last_commit"`
}

// CreateWikiPageOptions form for creating wiki
type CreateWikiPageOptions struct {
	// page title. leave empty to keep unchanged
	Title string `json:"title"`
	// content must be base64 encoded
	ContentBase64 string `json:"content_base64"`
	// optional commit message summarizing the change
	Message string `json:"message"`
}

// WikiCommitList commit/revision list
type WikiCommitList struct {
	// The list of wiki commits
	WikiCommits []*WikiCommit `json:"commits"`
	// The total count of commits
	Count int64 `json:"count"`
}
