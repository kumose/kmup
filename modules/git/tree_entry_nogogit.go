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

//go:build !gogit

package git

import "github.com/kumose/kmup/modules/log"

// Size returns the size of the entry
func (te *TreeEntry) Size() int64 {
	if te.IsDir() {
		return 0
	} else if te.sized {
		return te.size
	}

	wr, rd, cancel, err := te.ptree.repo.CatFileBatchCheck(te.ptree.repo.Ctx)
	if err != nil {
		log.Debug("error whilst reading size for %s in %s. Error: %v", te.ID.String(), te.ptree.repo.Path, err)
		return 0
	}
	defer cancel()
	_, err = wr.Write([]byte(te.ID.String() + "\n"))
	if err != nil {
		log.Debug("error whilst reading size for %s in %s. Error: %v", te.ID.String(), te.ptree.repo.Path, err)
		return 0
	}
	_, _, te.size, err = ReadBatchLine(rd)
	if err != nil {
		log.Debug("error whilst reading size for %s in %s. Error: %v", te.ID.String(), te.ptree.repo.Path, err)
		return 0
	}

	te.sized = true
	return te.size
}

// Blob returns the blob object the entry
func (te *TreeEntry) Blob() *Blob {
	return &Blob{
		ID:      te.ID,
		name:    te.Name(),
		size:    te.size,
		gotSize: te.sized,
		repo:    te.ptree.repo,
	}
}
