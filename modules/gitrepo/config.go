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
	"strings"

	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/globallock"
)

func GitConfigGet(ctx context.Context, repo Repository, key string) (string, error) {
	result, err := RunCmdString(ctx, repo, gitcmd.NewCommand("config", "--get").
		AddDynamicArguments(key))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}

func getRepoConfigLockKey(repoStoragePath string) string {
	return "repo-config:" + repoStoragePath
}

// GitConfigAdd add a git configuration key to a specific value for the given repository.
func GitConfigAdd(ctx context.Context, repo Repository, key, value string) error {
	return globallock.LockAndDo(ctx, getRepoConfigLockKey(repo.RelativePath()), func(ctx context.Context) error {
		_, err := RunCmdString(ctx, repo, gitcmd.NewCommand("config", "--add").
			AddDynamicArguments(key, value))
		return err
	})
}

// GitConfigSet updates a git configuration key to a specific value for the given repository.
// If the key does not exist, it will be created.
// If the key exists, it will be updated to the new value.
func GitConfigSet(ctx context.Context, repo Repository, key, value string) error {
	return globallock.LockAndDo(ctx, getRepoConfigLockKey(repo.RelativePath()), func(ctx context.Context) error {
		_, err := RunCmdString(ctx, repo, gitcmd.NewCommand("config").
			AddDynamicArguments(key, value))
		return err
	})
}
