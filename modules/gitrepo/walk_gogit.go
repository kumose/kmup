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

package gitrepo

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
)

// WalkReferences walks all the references from the repository
// refname is empty, ObjectTag or ObjectBranch. All other values should be treated as equivalent to empty.
func WalkReferences(ctx context.Context, repo Repository, walkfn func(sha1, refname string) error) (int, error) {
	gitRepo, closer, err := RepositoryFromContextOrOpen(ctx, repo)
	if err != nil {
		return 0, err
	}
	defer closer.Close()

	i := 0
	iter, err := gitRepo.GoGitRepo().References()
	if err != nil {
		return i, err
	}
	defer iter.Close()

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		err := walkfn(ref.Hash().String(), string(ref.Name()))
		i++
		return err
	})
	return i, err
}
