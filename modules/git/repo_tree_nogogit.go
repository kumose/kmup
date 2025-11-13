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

import (
	"io"
)

func (repo *Repository) getTree(id ObjectID) (*Tree, error) {
	wr, rd, cancel, err := repo.CatFileBatch(repo.Ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	_, _ = wr.Write([]byte(id.String() + "\n"))

	// ignore the SHA
	_, typ, size, err := ReadBatchLine(rd)
	if err != nil {
		return nil, err
	}

	switch typ {
	case "tag":
		resolvedID := id
		data, err := io.ReadAll(io.LimitReader(rd, size))
		if err != nil {
			return nil, err
		}
		tag, err := parseTagData(id.Type(), data)
		if err != nil {
			return nil, err
		}

		if _, err := wr.Write([]byte(tag.Object.String() + "\n")); err != nil {
			return nil, err
		}
		commit, err := repo.getCommitFromBatchReader(wr, rd, tag.Object)
		if err != nil {
			return nil, err
		}
		commit.Tree.ResolvedID = resolvedID
		return &commit.Tree, nil
	case "commit":
		commit, err := CommitFromReader(repo, id, io.LimitReader(rd, size))
		if err != nil {
			return nil, err
		}
		if _, err := rd.Discard(1); err != nil {
			return nil, err
		}
		commit.Tree.ResolvedID = commit.ID
		return &commit.Tree, nil
	case "tree":
		tree := NewTree(repo, id)
		tree.ResolvedID = id
		objectFormat, err := repo.GetObjectFormat()
		if err != nil {
			return nil, err
		}
		tree.entries, err = catBatchParseTreeEntries(objectFormat, tree, rd, size)
		if err != nil {
			return nil, err
		}
		tree.entriesParsed = true
		return tree, nil
	default:
		if err := DiscardFull(rd, size+1); err != nil {
			return nil, err
		}
		return nil, ErrNotExist{
			ID: id.String(),
		}
	}
}

// GetTree find the tree object in the repository.
func (repo *Repository) GetTree(idStr string) (*Tree, error) {
	objectFormat, err := repo.GetObjectFormat()
	if err != nil {
		return nil, err
	}
	if len(idStr) != objectFormat.FullLength() {
		res, err := repo.GetRefCommitID(idStr)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			idStr = res
		}
	}
	id, err := NewIDFromString(idStr)
	if err != nil {
		return nil, err
	}

	return repo.getTree(id)
}
