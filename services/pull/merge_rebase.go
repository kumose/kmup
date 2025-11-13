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
	"fmt"
	"strings"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/log"
)

// getRebaseAmendMessage composes the message to amend commits in rebase merge of a pull request.
func getRebaseAmendMessage(ctx *mergeContext, baseGitRepo *git.Repository) (message string, err error) {
	// Get existing commit message.
	commitMessage, _, err := gitcmd.NewCommand("show", "--format=%B", "-s").WithDir(ctx.tmpBasePath).RunStdString(ctx)
	if err != nil {
		return "", err
	}

	commitTitle, commitBody, _ := strings.Cut(commitMessage, "\n")
	extraVars := map[string]string{"CommitTitle": strings.TrimSpace(commitTitle), "CommitBody": strings.TrimSpace(commitBody)}

	message, body, err := getMergeMessage(ctx, baseGitRepo, ctx.pr, repo_model.MergeStyleRebase, extraVars)
	if err != nil || message == "" {
		return "", err
	}

	if len(body) > 0 {
		message = message + "\n\n" + body
	}
	return message, err
}

// Perform rebase merge without merge commit.
func doMergeRebaseFastForward(ctx *mergeContext) error {
	baseHeadSHA, err := git.GetFullCommitID(ctx, ctx.tmpBasePath, "HEAD")
	if err != nil {
		return fmt.Errorf("Failed to get full commit id for HEAD: %w", err)
	}

	cmd := gitcmd.NewCommand("merge", "--ff-only").AddDynamicArguments(stagingBranch)
	if err := runMergeCommand(ctx, repo_model.MergeStyleRebase, cmd); err != nil {
		log.Error("Unable to merge staging into base: %v", err)
		return err
	}

	// Check if anything actually changed before we amend the message, fast forward can skip commits.
	newMergeHeadSHA, err := git.GetFullCommitID(ctx, ctx.tmpBasePath, "HEAD")
	if err != nil {
		return fmt.Errorf("Failed to get full commit id for HEAD: %w", err)
	}
	if baseHeadSHA == newMergeHeadSHA {
		return nil
	}

	// Original repo to read template from.
	baseGitRepo, err := gitrepo.OpenRepository(ctx, ctx.pr.BaseRepo)
	if err != nil {
		log.Error("Unable to get Git repo for rebase: %v", err)
		return err
	}
	defer baseGitRepo.Close()

	// Amend last commit message based on template, if one exists
	newMessage, err := getRebaseAmendMessage(ctx, baseGitRepo)
	if err != nil {
		log.Error("Unable to get commit message for amend: %v", err)
		return err
	}

	if newMessage != "" {
		if err := gitcmd.NewCommand("commit", "--amend").
			AddOptionFormat("--message=%s", newMessage).
			WithDir(ctx.tmpBasePath).
			Run(ctx); err != nil {
			log.Error("Unable to amend commit message: %v", err)
			return err
		}
	}

	return nil
}

// Perform rebase merge with merge commit.
func doMergeRebaseMergeCommit(ctx *mergeContext, message string) error {
	cmd := gitcmd.NewCommand("merge").AddArguments("--no-ff", "--no-commit").AddDynamicArguments(stagingBranch)

	if err := runMergeCommand(ctx, repo_model.MergeStyleRebaseMerge, cmd); err != nil {
		log.Error("Unable to merge staging into base: %v", err)
		return err
	}
	if err := commitAndSignNoAuthor(ctx, message); err != nil {
		log.Error("Unable to make final commit: %v", err)
		return err
	}

	return nil
}

// doMergeStyleRebase rebases the tracking branch on the base branch as the current HEAD with or with a merge commit to the original pr branch
func doMergeStyleRebase(ctx *mergeContext, mergeStyle repo_model.MergeStyle, message string) error {
	if err := rebaseTrackingOnToBase(ctx, mergeStyle); err != nil {
		return err
	}

	// Checkout base branch again
	if err := ctx.PrepareGitCmd(gitcmd.NewCommand("checkout").AddDynamicArguments(baseBranch)).
		Run(ctx); err != nil {
		log.Error("git checkout base prior to merge post staging rebase %-v: %v\n%s\n%s", ctx.pr, err, ctx.outbuf.String(), ctx.errbuf.String())
		return fmt.Errorf("git checkout base prior to merge post staging rebase  %v: %w\n%s\n%s", ctx.pr, err, ctx.outbuf.String(), ctx.errbuf.String())
	}
	ctx.outbuf.Reset()
	ctx.errbuf.Reset()

	if mergeStyle == repo_model.MergeStyleRebase {
		return doMergeRebaseFastForward(ctx)
	}

	return doMergeRebaseMergeCommit(ctx, message)
}
