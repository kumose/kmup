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

package repo

import (
	git_model "github.com/kumose/kmup/models/git"
	access_model "github.com/kumose/kmup/models/perm/access"
	unit_model "github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/context"
	repo_service "github.com/kumose/kmup/services/repository"
)

type RecentBranchesPromptDataStruct struct {
	RecentlyPushedNewBranches []*git_model.RecentlyPushedNewBranch
}

func prepareRecentlyPushedNewBranches(ctx *context.Context) {
	if ctx.Doer == nil {
		return
	}
	if err := ctx.Repo.Repository.GetBaseRepo(ctx); err != nil {
		log.Error("GetBaseRepo: %v", err)
		return
	}

	opts := git_model.FindRecentlyPushedNewBranchesOptions{
		Repo:     ctx.Repo.Repository,
		BaseRepo: ctx.Repo.Repository,
	}
	if ctx.Repo.Repository.IsFork {
		opts.BaseRepo = ctx.Repo.Repository.BaseRepo
	}

	baseRepoPerm, err := access_model.GetUserRepoPermission(ctx, opts.BaseRepo, ctx.Doer)
	if err != nil {
		log.Error("GetUserRepoPermission: %v", err)
		return
	}
	if !opts.Repo.CanContentChange() || !opts.BaseRepo.CanContentChange() {
		return
	}
	if !opts.BaseRepo.UnitEnabled(ctx, unit_model.TypePullRequests) || !baseRepoPerm.CanRead(unit_model.TypePullRequests) {
		return
	}

	var finalBranches []*git_model.RecentlyPushedNewBranch
	branches, err := git_model.FindRecentlyPushedNewBranches(ctx, ctx.Doer, opts)
	if err != nil {
		log.Error("FindRecentlyPushedNewBranches failed: %v", err)
		return
	}

	for _, branch := range branches {
		divergingInfo, err := repo_service.GetBranchDivergingInfo(ctx,
			branch.BranchRepo, branch.BranchName, // "base" repo for diverging info
			opts.BaseRepo, opts.BaseRepo.DefaultBranch, // "head" repo for diverging info
		)
		if err != nil {
			log.Error("GetBranchDivergingInfo failed: %v", err)
			continue
		}
		branchRepoHasNewCommits := divergingInfo.BaseHasNewCommits
		baseRepoCommitsBehind := divergingInfo.HeadCommitsBehind
		if branchRepoHasNewCommits || baseRepoCommitsBehind > 0 {
			finalBranches = append(finalBranches, branch)
		}
	}
	if len(finalBranches) > 0 {
		ctx.Data["RecentBranchesPromptData"] = RecentBranchesPromptDataStruct{finalBranches}
	}
}
