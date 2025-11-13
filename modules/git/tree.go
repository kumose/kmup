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

package git

import (
	"bytes"
	"strings"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

type TreeCommon struct {
	ID         ObjectID
	ResolvedID ObjectID

	repo  *Repository
	ptree *Tree // parent tree
}

// NewTree create a new tree according the repository and tree id
func NewTree(repo *Repository, id ObjectID) *Tree {
	return &Tree{
		TreeCommon: TreeCommon{
			ID:   id,
			repo: repo,
		},
	}
}

// SubTree get a subtree by the sub dir path
func (t *Tree) SubTree(rpath string) (*Tree, error) {
	if len(rpath) == 0 {
		return t, nil
	}

	paths := strings.Split(rpath, "/")
	var (
		err error
		g   = t
		p   = t
		te  *TreeEntry
	)
	for _, name := range paths {
		te, err = p.GetTreeEntryByPath(name)
		if err != nil {
			return nil, err
		}

		g, err = t.repo.getTree(te.ID)
		if err != nil {
			return nil, err
		}
		g.ptree = p
		p = g
	}
	return g, nil
}

// LsTree checks if the given filenames are in the tree
func (repo *Repository) LsTree(ref string, filenames ...string) ([]string, error) {
	cmd := gitcmd.NewCommand("ls-tree", "-z", "--name-only").
		AddDashesAndList(append([]string{ref}, filenames...)...)

	res, _, err := cmd.WithDir(repo.Path).RunStdBytes(repo.Ctx)
	if err != nil {
		return nil, err
	}
	filelist := make([]string, 0, len(filenames))
	for line := range bytes.SplitSeq(res, []byte{'\000'}) {
		filelist = append(filelist, string(line))
	}

	return filelist, err
}

// GetTreePathLatestCommit returns the latest commit of a tree path
func (repo *Repository) GetTreePathLatestCommit(refName, treePath string) (*Commit, error) {
	stdout, _, err := gitcmd.NewCommand("rev-list", "-1").
		AddDynamicArguments(refName).AddDashesAndList(treePath).
		WithDir(repo.Path).
		RunStdString(repo.Ctx)
	if err != nil {
		return nil, err
	}
	return repo.GetCommit(strings.TrimSpace(stdout))
}
