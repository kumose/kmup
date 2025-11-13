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

package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/log"
	repo_module "github.com/kumose/kmup/modules/repository"
	"github.com/kumose/kmup/modules/setting"
	asymkey_service "github.com/kumose/kmup/services/asymkey"
)

// initRepoCommit temporarily changes with work directory.
func initRepoCommit(ctx context.Context, tmpPath string, repo *repo_model.Repository, u *user_model.User, defaultBranch string) (err error) {
	commitTimeStr := time.Now().Format(time.RFC3339)

	sig := u.NewGitSig()
	// Because this may call hooks we should pass in the environment
	env := append(os.Environ(),
		"GIT_AUTHOR_NAME="+sig.Name,
		"GIT_AUTHOR_EMAIL="+sig.Email,
		"GIT_AUTHOR_DATE="+commitTimeStr,
		"GIT_COMMITTER_DATE="+commitTimeStr,
	)
	committerName := sig.Name
	committerEmail := sig.Email

	if stdout, _, err := gitcmd.NewCommand("add", "--all").WithDir(tmpPath).RunStdString(ctx); err != nil {
		log.Error("git add --all failed: Stdout: %s\nError: %v", stdout, err)
		return fmt.Errorf("git add --all: %w", err)
	}

	cmd := gitcmd.NewCommand("commit", "--message=Initial commit").
		AddOptionFormat("--author='%s <%s>'", sig.Name, sig.Email)

	sign, key, signer, _ := asymkey_service.SignInitialCommit(ctx, tmpPath, u)
	if sign {
		if key.Format != "" {
			cmd.AddConfig("gpg.format", key.Format)
		}
		cmd.AddOptionFormat("-S%s", key.KeyID)

		if repo.GetTrustModel() == repo_model.CommitterTrustModel || repo.GetTrustModel() == repo_model.CollaboratorCommitterTrustModel {
			// need to set the committer to the KeyID owner
			committerName = signer.Name
			committerEmail = signer.Email
		}
	} else {
		cmd.AddArguments("--no-gpg-sign")
	}

	env = append(env,
		"GIT_COMMITTER_NAME="+committerName,
		"GIT_COMMITTER_EMAIL="+committerEmail,
	)

	if stdout, _, err := cmd.WithDir(tmpPath).WithEnv(env).RunStdString(ctx); err != nil {
		log.Error("Failed to commit: %v: Stdout: %s\nError: %v", cmd.LogString(), stdout, err)
		return fmt.Errorf("git commit: %w", err)
	}

	if len(defaultBranch) == 0 {
		defaultBranch = setting.Repository.DefaultBranch
	}

	if stdout, _, err := gitcmd.NewCommand("push", "origin").
		AddDynamicArguments("HEAD:" + defaultBranch).
		WithDir(tmpPath).
		WithEnv(repo_module.InternalPushingEnvironment(u, repo)).
		RunStdString(ctx); err != nil {
		log.Error("Failed to push back to HEAD: Stdout: %s\nError: %v", stdout, err)
		return fmt.Errorf("git push: %w", err)
	}

	return nil
}
