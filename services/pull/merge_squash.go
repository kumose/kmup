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

	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

// doMergeStyleSquash gets a commit author signature for squash commits
func getAuthorSignatureSquash(ctx *mergeContext) (*git.Signature, error) {
	if err := ctx.pr.Issue.LoadPoster(ctx); err != nil {
		log.Error("%-v Issue[%d].LoadPoster: %v", ctx.pr, ctx.pr.Issue.ID, err)
		return nil, err
	}

	// Try to get an signature from the same user in one of the commits, as the
	// poster email might be private or commits might have a different signature
	// than the primary email address of the poster.
	gitRepo, err := git.OpenRepository(ctx, ctx.tmpBasePath)
	if err != nil {
		log.Error("%-v Unable to open base repository: %v", ctx.pr, err)
		return nil, err
	}
	defer gitRepo.Close()

	commits, err := gitRepo.CommitsBetweenIDs(trackingBranch, "HEAD")
	if err != nil {
		log.Error("%-v Unable to get commits between: %s %s: %v", ctx.pr, "HEAD", trackingBranch, err)
		return nil, err
	}

	uniqueEmails := make(container.Set[string])
	for _, commit := range commits {
		if commit.Author != nil && uniqueEmails.Add(commit.Author.Email) {
			commitUser, _ := user_model.GetUserByEmail(ctx, commit.Author.Email)
			if commitUser != nil && commitUser.ID == ctx.pr.Issue.Poster.ID {
				return commit.Author, nil
			}
		}
	}

	return ctx.pr.Issue.Poster.NewGitSig(), nil
}

// doMergeStyleSquash squashes the tracking branch on the current HEAD (=base)
func doMergeStyleSquash(ctx *mergeContext, message string) error {
	sig, err := getAuthorSignatureSquash(ctx)
	if err != nil {
		return fmt.Errorf("getAuthorSignatureSquash: %w", err)
	}

	cmdMerge := gitcmd.NewCommand("merge", "--squash").AddDynamicArguments(trackingBranch)
	if err := runMergeCommand(ctx, repo_model.MergeStyleSquash, cmdMerge); err != nil {
		log.Error("%-v Unable to merge --squash tracking into base: %v", ctx.pr, err)
		return err
	}

	if setting.Repository.PullRequest.AddCoCommitterTrailers && ctx.committer.String() != sig.String() {
		// add trailer
		message = AddCommitMessageTailer(message, "Co-authored-by", sig.String())
		message = AddCommitMessageTailer(message, "Co-committed-by", sig.String()) // FIXME: this one should be removed, it is not really used or widely used
	}
	cmdCommit := gitcmd.NewCommand("commit").
		AddOptionFormat("--author='%s <%s>'", sig.Name, sig.Email).
		AddOptionFormat("--message=%s", message)
	if ctx.signKey == nil {
		cmdCommit.AddArguments("--no-gpg-sign")
	} else {
		if ctx.signKey.Format != "" {
			cmdCommit.AddConfig("gpg.format", ctx.signKey.Format)
		}
		cmdCommit.AddOptionFormat("-S%s", ctx.signKey.KeyID)
	}
	if err := ctx.PrepareGitCmd(cmdCommit).Run(ctx); err != nil {
		log.Error("git commit %-v: %v\n%s\n%s", ctx.pr, err, ctx.outbuf.String(), ctx.errbuf.String())
		return fmt.Errorf("git commit [%s:%s -> %s:%s]: %w\n%s\n%s", ctx.pr.HeadRepo.FullName(), ctx.pr.HeadBranch, ctx.pr.BaseRepo.FullName(), ctx.pr.BaseBranch, err, ctx.outbuf.String(), ctx.errbuf.String())
	}
	ctx.outbuf.Reset()
	ctx.errbuf.Reset()
	return nil
}
