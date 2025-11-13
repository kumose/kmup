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

//go:build gogit

package git

import (
	"os"
	"path/filepath"

	kmuplog "github.com/kumose/kmup/modules/log"

	commitgraph "github.com/go-git/go-git/v5/plumbing/format/commitgraph/v2"
	cgobject "github.com/go-git/go-git/v5/plumbing/object/commitgraph"
)

// CommitNodeIndex returns the index for walking commit graph
func (r *Repository) CommitNodeIndex() (cgobject.CommitNodeIndex, *os.File) {
	indexPath := filepath.Join(r.Path, "objects", "info", "commit-graph")

	file, err := os.Open(indexPath)
	if err == nil {
		var index commitgraph.Index
		index, err = commitgraph.OpenFileIndex(file)
		if err == nil {
			return cgobject.NewGraphCommitNodeIndex(index, r.gogitRepo.Storer), file
		}
	}

	if !os.IsNotExist(err) {
		kmuplog.Warn("Unable to read commit-graph for %s: %v", r.Path, err)
	}

	return cgobject.NewObjectCommitNodeIndex(r.gogitRepo.Storer), nil
}
