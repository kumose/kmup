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
	"io"
	"time"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/gitcmd"
	giturl "github.com/kumose/kmup/modules/git/url"
	"github.com/kumose/kmup/modules/globallock"
	"github.com/kumose/kmup/modules/util"
)

type RemoteOption string

const (
	RemoteOptionMirrorPush  RemoteOption = "--mirror=push"
	RemoteOptionMirrorFetch RemoteOption = "--mirror=fetch"
)

func GitRemoteAdd(ctx context.Context, repo Repository, remoteName, remoteURL string, options ...RemoteOption) error {
	return globallock.LockAndDo(ctx, getRepoConfigLockKey(repo.RelativePath()), func(ctx context.Context) error {
		cmd := gitcmd.NewCommand("remote", "add")
		if len(options) > 0 {
			switch options[0] {
			case RemoteOptionMirrorPush:
				cmd.AddArguments("--mirror=push")
			case RemoteOptionMirrorFetch:
				cmd.AddArguments("--mirror=fetch")
			default:
				return errors.New("unknown remote option: " + string(options[0]))
			}
		}
		_, err := RunCmdString(ctx, repo, cmd.AddDynamicArguments(remoteName, remoteURL))
		return err
	})
}

func GitRemoteRemove(ctx context.Context, repo Repository, remoteName string) error {
	return globallock.LockAndDo(ctx, getRepoConfigLockKey(repo.RelativePath()), func(ctx context.Context) error {
		cmd := gitcmd.NewCommand("remote", "rm").AddDynamicArguments(remoteName)
		_, err := RunCmdString(ctx, repo, cmd)
		return err
	})
}

// GitRemoteGetURL returns the url of a specific remote of the repository.
func GitRemoteGetURL(ctx context.Context, repo Repository, remoteName string) (*giturl.GitURL, error) {
	addr, err := git.GetRemoteAddress(ctx, repoPath(repo), remoteName)
	if err != nil {
		return nil, err
	}
	if addr == "" {
		return nil, util.NewNotExistErrorf("remote '%s' does not exist", remoteName)
	}
	return giturl.ParseGitURL(addr)
}

// GitRemotePrune prunes the remote branches that no longer exist in the remote repository.
func GitRemotePrune(ctx context.Context, repo Repository, remoteName string, timeout time.Duration, stdout, stderr io.Writer) error {
	return RunCmd(ctx, repo, gitcmd.NewCommand("remote", "prune").
		AddDynamicArguments(remoteName).
		WithTimeout(timeout).
		WithStdout(stdout).
		WithStderr(stderr))
}

// GitRemoteUpdatePrune updates the remote branches and prunes the ones that no longer exist in the remote repository.
func GitRemoteUpdatePrune(ctx context.Context, repo Repository, remoteName string, timeout time.Duration, stdout, stderr io.Writer) error {
	return RunCmd(ctx, repo, gitcmd.NewCommand("remote", "update", "--prune").
		AddDynamicArguments(remoteName).
		WithTimeout(timeout).
		WithStdout(stdout).
		WithStderr(stderr))
}
