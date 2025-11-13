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
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GetRefsFiltered returns all references of the repository that matches patterm exactly or starting with.
func (repo *Repository) GetRefsFiltered(pattern string) ([]*Reference, error) {
	r, err := git.PlainOpen(repo.Path)
	if err != nil {
		return nil, err
	}

	refsIter, err := r.References()
	if err != nil {
		return nil, err
	}
	refs := make([]*Reference, 0)
	if err = refsIter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() != plumbing.HEAD && !ref.Name().IsRemote() &&
			(pattern == "" || strings.HasPrefix(ref.Name().String(), pattern)) {
			refType := string(ObjectCommit)
			if ref.Name().IsTag() {
				// tags can be of type `commit` (lightweight) or `tag` (annotated)
				if tagType, _ := repo.GetTagType(ParseGogitHash(ref.Hash())); err == nil {
					refType = tagType
				}
			}
			r := &Reference{
				Name:   ref.Name().String(),
				Object: ParseGogitHash(ref.Hash()),
				Type:   refType,
				repo:   repo,
			}
			refs = append(refs, r)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return refs, nil
}
