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
	"os"
	"strconv"
	"strings"

	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
)

// env keys for git hooks need
const (
	EnvRepoName     = "KMUP_REPO_NAME"
	EnvRepoUsername = "KMUP_REPO_USER_NAME"
	EnvRepoID       = "KMUP_REPO_ID"
	EnvRepoIsWiki   = "KMUP_REPO_IS_WIKI"
	EnvPusherName   = "KMUP_PUSHER_NAME"
	EnvPusherEmail  = "KMUP_PUSHER_EMAIL"
	EnvPusherID     = "KMUP_PUSHER_ID"
	EnvKeyID        = "KMUP_KEY_ID" // public key ID
	EnvDeployKeyID  = "KMUP_DEPLOY_KEY_ID"
	EnvPRID         = "KMUP_PR_ID"
	EnvPushTrigger  = "KMUP_PUSH_TRIGGER"
	EnvIsInternal   = "KMUP_INTERNAL_PUSH"
	EnvAppURL       = "KMUP_ROOT_URL"
	EnvActionPerm   = "KMUP_ACTION_PERM"
)

type PushTrigger string

const (
	PushTriggerPRMergeToBase    PushTrigger = "pr-merge-to-base"
	PushTriggerPRUpdateWithBase PushTrigger = "pr-update-with-base"
)

// InternalPushingEnvironment returns an os environment to switch off hooks on push
// It is recommended to avoid using this unless you are pushing within a transaction
// or if you absolutely are sure that post-receive and pre-receive will do nothing
// We provide the full pushing-environment for other hook providers
func InternalPushingEnvironment(doer *user_model.User, repo *repo_model.Repository) []string {
	return append(PushingEnvironment(doer, repo),
		EnvIsInternal+"=true",
	)
}

// PushingEnvironment returns an os environment to allow hooks to work on push
func PushingEnvironment(doer *user_model.User, repo *repo_model.Repository) []string {
	return FullPushingEnvironment(doer, doer, repo, repo.Name, 0)
}

// FullPushingEnvironment returns an os environment to allow hooks to work on push
func FullPushingEnvironment(author, committer *user_model.User, repo *repo_model.Repository, repoName string, prID int64) []string {
	isWiki := "false"
	if strings.HasSuffix(repoName, ".wiki") {
		isWiki = "true"
	}

	authorSig := author.NewGitSig()
	committerSig := committer.NewGitSig()

	environ := append(os.Environ(),
		"GIT_AUTHOR_NAME="+authorSig.Name,
		"GIT_AUTHOR_EMAIL="+authorSig.Email,
		"GIT_COMMITTER_NAME="+committerSig.Name,
		"GIT_COMMITTER_EMAIL="+committerSig.Email,
		EnvRepoName+"="+repoName,
		EnvRepoUsername+"="+repo.OwnerName,
		EnvRepoIsWiki+"="+isWiki,
		EnvPusherName+"="+committer.Name,
		EnvPusherID+"="+strconv.FormatInt(committer.ID, 10),
		EnvRepoID+"="+strconv.FormatInt(repo.ID, 10),
		EnvPRID+"="+strconv.FormatInt(prID, 10),
		EnvAppURL+"="+setting.AppURL,
		"SSH_ORIGINAL_COMMAND=kmup-internal",
	)

	if !committer.KeepEmailPrivate {
		environ = append(environ, EnvPusherEmail+"="+committer.Email)
	}

	return environ
}
