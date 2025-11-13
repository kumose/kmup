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
	"errors"

	"github.com/kumose/kmup/modules/git/gitcmd"

	"github.com/go-git/go-git/v5/plumbing"
)

func (repo *Repository) getTree(id ObjectID) (*Tree, error) {
	gogitTree, err := repo.gogitRepo.TreeObject(plumbing.Hash(id.RawValue()))
	if err != nil {
		if errors.Is(err, plumbing.ErrObjectNotFound) {
			return nil, ErrNotExist{
				ID: id.String(),
			}
		}
		return nil, err
	}

	tree := NewTree(repo, id)
	tree.resolvedGogitTreeObject = gogitTree
	return tree, nil
}

// GetTree find the tree object in the repository.
func (repo *Repository) GetTree(idStr string) (*Tree, error) {
	objectFormat, err := repo.GetObjectFormat()
	if err != nil {
		return nil, err
	}

	if len(idStr) != objectFormat.FullLength() {
		res, _, err := gitcmd.NewCommand("rev-parse", "--verify").
			AddDynamicArguments(idStr).
			WithDir(repo.Path).
			RunStdString(repo.Ctx)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			idStr = res[:len(res)-1]
		}
	}
	id, err := NewIDFromString(idStr)
	if err != nil {
		return nil, err
	}
	resolvedID := id
	commitObject, err := repo.gogitRepo.CommitObject(plumbing.Hash(id.RawValue()))
	if err == nil {
		id = ParseGogitHash(commitObject.TreeHash)
	}
	treeObject, err := repo.getTree(id)
	if err != nil {
		return nil, err
	}
	treeObject.ResolvedID = resolvedID
	return treeObject, nil
}
