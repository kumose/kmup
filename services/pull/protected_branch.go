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

package pull

import (
	"context"

	git_model "github.com/kumose/kmup/models/git"
	repo_model "github.com/kumose/kmup/models/repo"
)

func CreateOrUpdateProtectedBranch(ctx context.Context, repo *repo_model.Repository,
	protectBranch *git_model.ProtectedBranch, whitelistOptions git_model.WhitelistOptions,
) error {
	err := git_model.UpdateProtectBranch(ctx, repo, protectBranch, whitelistOptions)
	if err != nil {
		return err
	}

	isPlainRule := !git_model.IsRuleNameSpecial(protectBranch.RuleName)
	var isBranchExist bool
	if isPlainRule {
		isBranchExist, _ = git_model.IsBranchExist(ctx, repo.ID, protectBranch.RuleName)
	}

	if isBranchExist {
		if err := CheckPRsForBaseBranch(ctx, repo, protectBranch.RuleName); err != nil {
			return err
		}
	} else {
		if !isPlainRule {
			// FIXME: since we only need to recheck files protected rules, we could improve this
			matchedBranches, err := git_model.FindAllMatchedBranches(ctx, repo.ID, protectBranch.RuleName)
			if err != nil {
				return err
			}
			for _, branchName := range matchedBranches {
				if err = CheckPRsForBaseBranch(ctx, repo, branchName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
