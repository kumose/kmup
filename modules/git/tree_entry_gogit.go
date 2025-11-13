// Copyright 2015 The Gogs Authors. All rights reserved.
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
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// gogitFileModeToEntryMode converts go-git filemode to EntryMode
func gogitFileModeToEntryMode(mode filemode.FileMode) EntryMode {
	return EntryMode(mode)
}

func entryModeToGogitFileMode(mode EntryMode) filemode.FileMode {
	return filemode.FileMode(mode)
}

func (te *TreeEntry) toGogitTreeEntry() *object.TreeEntry {
	return &object.TreeEntry{
		Name: te.name,
		Mode: entryModeToGogitFileMode(te.entryMode),
		Hash: plumbing.Hash(te.ID.RawValue()),
	}
}

// Size returns the size of the entry
func (te *TreeEntry) Size() int64 {
	if te.IsDir() {
		return 0
	} else if te.sized {
		return te.size
	}

	ptreeGogitTree, err := te.ptree.gogitTreeObject()
	if err != nil {
		return 0
	}
	file, err := ptreeGogitTree.TreeEntryFile(te.toGogitTreeEntry())
	if err != nil {
		return 0
	}

	te.sized = true
	te.size = file.Size
	return te.size
}

// Blob returns the blob object the entry
func (te *TreeEntry) Blob() *Blob {
	encodedObj, err := te.ptree.repo.gogitRepo.Storer.EncodedObject(plumbing.AnyObject, te.toGogitTreeEntry().Hash)
	if err != nil {
		return nil
	}

	return &Blob{
		ID:              te.ID,
		gogitEncodedObj: encodedObj,
		name:            te.Name(),
	}
}
