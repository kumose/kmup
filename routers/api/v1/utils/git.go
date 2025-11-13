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

package utils

import (
	"errors"

	git_model "github.com/kumose/kmup/models/git"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/reqctx"
	"github.com/kumose/kmup/services/context"
)

type RefCommit struct {
	InputRef string
	RefName  git.RefName
	Commit   *git.Commit
	CommitID string
}

// ResolveRefCommit resolve ref to a commit if exist
func ResolveRefCommit(ctx reqctx.RequestContext, repo *repo_model.Repository, inputRef string, minCommitIDLen ...int) (_ *RefCommit, err error) {
	gitRepo, err := gitrepo.RepositoryFromRequestContextOrOpen(ctx, repo)
	if err != nil {
		return nil, err
	}
	refCommit := RefCommit{InputRef: inputRef}
	if exist, _ := git_model.IsBranchExist(ctx, repo.ID, inputRef); exist {
		refCommit.RefName = git.RefNameFromBranch(inputRef)
	} else if gitrepo.IsTagExist(ctx, repo, inputRef) {
		refCommit.RefName = git.RefNameFromTag(inputRef)
	} else if git.IsStringLikelyCommitID(git.ObjectFormatFromName(repo.ObjectFormatName), inputRef, minCommitIDLen...) {
		refCommit.RefName = git.RefNameFromCommit(inputRef)
	}
	if refCommit.RefName == "" {
		return nil, git.ErrNotExist{ID: inputRef}
	}
	if refCommit.Commit, err = gitRepo.GetCommit(refCommit.RefName.String()); err != nil {
		return nil, err
	}
	refCommit.CommitID = refCommit.Commit.ID.String()
	return &refCommit, nil
}

func NewRefCommit(refName git.RefName, commit *git.Commit) *RefCommit {
	return &RefCommit{InputRef: refName.ShortName(), RefName: refName, Commit: commit, CommitID: commit.ID.String()}
}

// GetGitRefs return git references based on filter
func GetGitRefs(ctx *context.APIContext, filter string) ([]*git.Reference, string, error) {
	if ctx.Repo.GitRepo == nil {
		return nil, "", errors.New("no open git repo found in context")
	}
	if len(filter) > 0 {
		filter = "refs/" + filter
	}
	refs, err := ctx.Repo.GitRepo.GetRefsFiltered(filter)
	return refs, "GetRefsFiltered", err
}
