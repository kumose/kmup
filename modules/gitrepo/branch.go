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
	"errors"
	"strings"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/gitcmd"
)

// GetBranchesByPath returns a branch by its path
// if limit = 0 it will not limit
func GetBranchesByPath(ctx context.Context, repo Repository, skip, limit int) ([]string, int, error) {
	gitRepo, err := OpenRepository(ctx, repo)
	if err != nil {
		return nil, 0, err
	}
	defer gitRepo.Close()

	return gitRepo.GetBranchNames(skip, limit)
}

func GetBranchCommitID(ctx context.Context, repo Repository, branch string) (string, error) {
	gitRepo, err := OpenRepository(ctx, repo)
	if err != nil {
		return "", err
	}
	defer gitRepo.Close()

	return gitRepo.GetBranchCommitID(branch)
}

// SetDefaultBranch sets default branch of repository.
func SetDefaultBranch(ctx context.Context, repo Repository, name string) error {
	_, err := RunCmdString(ctx, repo, gitcmd.NewCommand("symbolic-ref", "HEAD").
		AddDynamicArguments(git.BranchPrefix+name))
	return err
}

// GetDefaultBranch gets default branch of repository.
func GetDefaultBranch(ctx context.Context, repo Repository) (string, error) {
	stdout, err := RunCmdString(ctx, repo, gitcmd.NewCommand("symbolic-ref", "HEAD"))
	if err != nil {
		return "", err
	}
	stdout = strings.TrimSpace(stdout)
	if !strings.HasPrefix(stdout, git.BranchPrefix) {
		return "", errors.New("the HEAD is not a branch: " + stdout)
	}
	return strings.TrimPrefix(stdout, git.BranchPrefix), nil
}

// IsReferenceExist returns true if given reference exists in the repository.
func IsReferenceExist(ctx context.Context, repo Repository, name string) bool {
	_, err := RunCmdString(ctx, repo, gitcmd.NewCommand("show-ref", "--verify").AddDashesAndList(name))
	return err == nil
}

// IsBranchExist returns true if given branch exists in the repository.
func IsBranchExist(ctx context.Context, repo Repository, name string) bool {
	return IsReferenceExist(ctx, repo, git.BranchPrefix+name)
}

// DeleteBranch delete a branch by name on repository.
func DeleteBranch(ctx context.Context, repo Repository, name string, force bool) error {
	cmd := gitcmd.NewCommand("branch")

	if force {
		cmd.AddArguments("-D")
	} else {
		cmd.AddArguments("-d")
	}

	cmd.AddDashesAndList(name)
	_, err := RunCmdString(ctx, repo, cmd)
	return err
}

// CreateBranch create a new branch
func CreateBranch(ctx context.Context, repo Repository, branch, oldbranchOrCommit string) error {
	cmd := gitcmd.NewCommand("branch")
	cmd.AddDashesAndList(branch, oldbranchOrCommit)

	_, err := RunCmdString(ctx, repo, cmd)
	return err
}

// RenameBranch rename a branch
func RenameBranch(ctx context.Context, repo Repository, from, to string) error {
	_, err := RunCmdString(ctx, repo, gitcmd.NewCommand("branch", "-m").AddDynamicArguments(from, to))
	return err
}
