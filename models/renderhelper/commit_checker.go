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

package renderhelper

import (
	"context"
	"io"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/log"
)

type commitChecker struct {
	ctx           context.Context
	commitCache   map[string]bool
	gitRepoFacade gitrepo.Repository

	gitRepo       *git.Repository
	gitRepoCloser io.Closer
}

func newCommitChecker(ctx context.Context, gitRepo gitrepo.Repository) *commitChecker {
	return &commitChecker{ctx: ctx, commitCache: make(map[string]bool), gitRepoFacade: gitRepo}
}

func (c *commitChecker) Close() error {
	if c != nil && c.gitRepoCloser != nil {
		return c.gitRepoCloser.Close()
	}
	return nil
}

func (c *commitChecker) IsCommitIDExisting(commitID string) bool {
	exist, inCache := c.commitCache[commitID]
	if inCache {
		return exist
	}

	if c.gitRepo == nil {
		r, closer, err := gitrepo.RepositoryFromContextOrOpen(c.ctx, c.gitRepoFacade)
		if err != nil {
			log.Error("unable to open repository: %s Error: %v", gitrepo.RepoGitURL(c.gitRepoFacade), err)
			return false
		}
		c.gitRepo, c.gitRepoCloser = r, closer
	}

	exist = c.gitRepo.IsReferenceExist(commitID) // Don't use IsObjectExist since it doesn't support short hashes with gogit edition.
	c.commitCache[commitID] = exist
	return exist
}
