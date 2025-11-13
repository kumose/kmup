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

package gitrepo

import (
	"context"

	"github.com/kumose/kmup/modules/git"
)

// CloneExternalRepo clones an external repository to the managed repository.
func CloneExternalRepo(ctx context.Context, fromRemoteURL string, toRepo Repository, opts git.CloneRepoOptions) error {
	return git.Clone(ctx, fromRemoteURL, repoPath(toRepo), opts)
}

// CloneRepoToLocal clones a managed repository to a local path.
func CloneRepoToLocal(ctx context.Context, fromRepo Repository, toLocalPath string, opts git.CloneRepoOptions) error {
	return git.Clone(ctx, repoPath(fromRepo), toLocalPath, opts)
}
