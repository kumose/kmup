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

package doctor

import (
	"context"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/log"
)

func synchronizeRepoHeads(ctx context.Context, logger log.Logger, autofix bool) error {
	numRepos := 0
	numHeadsBroken := 0
	numDefaultBranchesBroken := 0
	numReposUpdated := 0
	err := iterateRepositories(ctx, func(repo *repo_model.Repository) error {
		numRepos++
		_, defaultBranchErr := gitrepo.RunCmdString(ctx, repo,
			gitcmd.NewCommand("rev-parse").AddDashesAndList(repo.DefaultBranch))

		head, headErr := gitrepo.RunCmdString(ctx, repo,
			gitcmd.NewCommand("symbolic-ref", "--short", "HEAD"))

		// what we expect: default branch is valid, and HEAD points to it
		if headErr == nil && defaultBranchErr == nil && head == repo.DefaultBranch {
			return nil
		}

		if headErr != nil {
			numHeadsBroken++
		}
		if defaultBranchErr != nil {
			numDefaultBranchesBroken++
		}

		// if default branch is broken, let the user fix that in the UI
		if defaultBranchErr != nil {
			logger.Warn("Default branch for %s/%s doesn't point to a valid commit", repo.OwnerName, repo.Name)
			return nil
		}

		// if we're not autofixing, that's all we can do
		if !autofix {
			return nil
		}

		// otherwise, let's try fixing HEAD
		err := gitrepo.RunCmd(ctx, repo, gitcmd.NewCommand("symbolic-ref").AddDashesAndList("HEAD", git.BranchPrefix+repo.DefaultBranch))
		if err != nil {
			logger.Warn("Failed to fix HEAD for %s/%s: %v", repo.OwnerName, repo.Name, err)
			return nil
		}
		numReposUpdated++
		return nil
	})
	if err != nil {
		logger.Critical("Error when fixing repo HEADs: %v", err)
	}

	if autofix {
		logger.Info("Out of %d repos, HEADs for %d are now fixed and HEADS for %d are still broken", numRepos, numReposUpdated, numDefaultBranchesBroken+numHeadsBroken-numReposUpdated)
	} else {
		if numHeadsBroken == 0 && numDefaultBranchesBroken == 0 {
			logger.Info("All %d repos have their HEADs in the correct state", numRepos)
		} else {
			if numHeadsBroken == 0 && numDefaultBranchesBroken != 0 {
				logger.Critical("Default branches are broken for %d/%d repos", numDefaultBranchesBroken, numRepos)
			} else if numHeadsBroken != 0 && numDefaultBranchesBroken == 0 {
				logger.Warn("HEADs are broken for %d/%d repos", numHeadsBroken, numRepos)
			} else {
				logger.Critical("Out of %d repos, HEADS are broken for %d and default branches are broken for %d", numRepos, numHeadsBroken, numDefaultBranchesBroken)
			}
		}
	}

	return err
}

func init() {
	Register(&Check{
		Title:     "Synchronize repo HEADs",
		Name:      "synchronize-repo-heads",
		IsDefault: true,
		Run:       synchronizeRepoHeads,
		Priority:  7,
	})
}
