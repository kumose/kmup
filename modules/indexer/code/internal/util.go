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

import (
	"strings"

	"github.com/kumose/kmup/modules/indexer/internal"
	"github.com/kumose/kmup/modules/log"
)

const filenameMatchNumberOfLines = 7 // Copied from GitHub search

func FilenameIndexerID(repoID int64, filename string) string {
	return internal.Base36(repoID) + "_" + filename
}

func ParseIndexerID(indexerID string) (int64, string) {
	index := strings.IndexByte(indexerID, '_')
	if index == -1 {
		log.Error("Unexpected ID in repo indexer: %s", indexerID)
	}
	repoID, _ := internal.ParseBase36(indexerID[:index])
	return repoID, indexerID[index+1:]
}

func FilenameOfIndexerID(indexerID string) string {
	index := strings.IndexByte(indexerID, '_')
	if index == -1 {
		log.Error("Unexpected ID in repo indexer: %s", indexerID)
	}
	return indexerID[index+1:]
}

// FilenameMatchIndexPos returns the boundaries of its first seven lines.
func FilenameMatchIndexPos(content string) (int, int) {
	count := 1
	for i, c := range content {
		if c == '\n' {
			count++
			if count == filenameMatchNumberOfLines {
				return 0, i
			}
		}
	}
	return 0, len(content)
}
