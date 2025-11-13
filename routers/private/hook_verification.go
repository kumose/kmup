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

package private

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/log"
	asymkey_service "github.com/kumose/kmup/services/asymkey"
)

// This file contains commit verification functions for refs passed across in hooks

func verifyCommits(oldCommitID, newCommitID string, repo *git.Repository, env []string) error {
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		log.Error("Unable to create os.Pipe for %s", repo.Path)
		return err
	}
	defer func() {
		_ = stdoutReader.Close()
		_ = stdoutWriter.Close()
	}()

	var command *gitcmd.Command
	objectFormat, _ := repo.GetObjectFormat()
	if oldCommitID == objectFormat.EmptyObjectID().String() {
		// When creating a new branch, the oldCommitID is empty, by using "newCommitID --not --all":
		// List commits that are reachable by following the newCommitID, exclude "all" existing heads/tags commits
		// So, it only lists the new commits received, doesn't list the commits already present in the receiving repository
		command = gitcmd.NewCommand("rev-list").AddDynamicArguments(newCommitID).AddArguments("--not", "--all")
	} else {
		command = gitcmd.NewCommand("rev-list").AddDynamicArguments(oldCommitID + "..." + newCommitID)
	}
	// This is safe as force pushes are already forbidden
	err = command.WithEnv(env).
		WithDir(repo.Path).
		WithStdout(stdoutWriter).
		WithPipelineFunc(func(ctx context.Context, cancel context.CancelFunc) error {
			_ = stdoutWriter.Close()
			err := readAndVerifyCommitsFromShaReader(stdoutReader, repo, env)
			if err != nil {
				log.Error("readAndVerifyCommitsFromShaReader failed: %v", err)
				cancel()
			}
			_ = stdoutReader.Close()
			return err
		}).
		Run(repo.Ctx)
	if err != nil && !isErrUnverifiedCommit(err) {
		log.Error("Unable to check commits from %s to %s in %s: %v", oldCommitID, newCommitID, repo.Path, err)
	}
	return err
}

func readAndVerifyCommitsFromShaReader(input io.ReadCloser, repo *git.Repository, env []string) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		err := readAndVerifyCommit(line, repo, env)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func readAndVerifyCommit(sha string, repo *git.Repository, env []string) error {
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		log.Error("Unable to create pipe for %s: %v", repo.Path, err)
		return err
	}
	defer func() {
		_ = stdoutReader.Close()
		_ = stdoutWriter.Close()
	}()

	commitID := git.MustIDFromString(sha)

	return gitcmd.NewCommand("cat-file", "commit").AddDynamicArguments(sha).
		WithEnv(env).
		WithDir(repo.Path).
		WithStdout(stdoutWriter).
		WithPipelineFunc(func(ctx context.Context, cancel context.CancelFunc) error {
			_ = stdoutWriter.Close()
			commit, err := git.CommitFromReader(repo, commitID, stdoutReader)
			if err != nil {
				return err
			}
			verification := asymkey_service.ParseCommitWithSignature(ctx, commit)
			if !verification.Verified {
				cancel()
				return &errUnverifiedCommit{
					commit.ID.String(),
				}
			}
			return nil
		}).
		Run(repo.Ctx)
}

type errUnverifiedCommit struct {
	sha string
}

func (e *errUnverifiedCommit) Error() string {
	return "Unverified commit: " + e.sha
}

func isErrUnverifiedCommit(err error) bool {
	_, ok := err.(*errUnverifiedCommit)
	return ok
}
