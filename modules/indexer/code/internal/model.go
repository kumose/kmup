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

package internal

import "github.com/kumose/kmup/modules/timeutil"

type FileUpdate struct {
	Filename string
	BlobSha  string
	Size     int64
	Sized    bool
}

// RepoChanges changes (file additions/updates/removals) to a repo
type RepoChanges struct {
	Updates          []FileUpdate
	RemovedFilenames []string
}

// IndexerData represents data stored in the code indexer
type IndexerData struct {
	RepoID int64
}

// SearchResult result of performing a search in a repo
type SearchResult struct {
	RepoID      int64
	StartIndex  int
	EndIndex    int
	Filename    string
	Content     string
	CommitID    string
	UpdatedUnix timeutil.TimeStamp
	Language    string
	Color       string
}

// SearchResultLanguages result of top languages count in search results
type SearchResultLanguages struct {
	Language string
	Color    string
	Count    int
}
