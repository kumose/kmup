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

package cron

import (
	"context"
	"time"

	activities_model "github.com/kumose/kmup/models/activities"
	"github.com/kumose/kmup/models/system"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/git/gitcmd"
	issue_indexer "github.com/kumose/kmup/modules/indexer/issues"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/updatechecker"
	asymkey_service "github.com/kumose/kmup/services/asymkey"
	repo_service "github.com/kumose/kmup/services/repository"
	archiver_service "github.com/kumose/kmup/services/repository/archiver"
	user_service "github.com/kumose/kmup/services/user"
)

func registerDeleteInactiveUsers() {
	RegisterTaskFatal("delete_inactive_accounts", &OlderThanConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@annually",
		},
		OlderThan: time.Minute * time.Duration(setting.Service.ActiveCodeLives),
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		olderThanConfig := config.(*OlderThanConfig)
		return user_service.DeleteInactiveUsers(ctx, olderThanConfig.OlderThan)
	})
}

func registerDeleteRepositoryArchives() {
	RegisterTaskFatal("delete_repo_archives", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@annually",
	}, func(ctx context.Context, _ *user_model.User, _ Config) error {
		return archiver_service.DeleteRepositoryArchives(ctx)
	})
}

func registerGarbageCollectRepositories() {
	type RepoHealthCheckConfig struct {
		BaseConfig
		Timeout time.Duration
		Args    []string `delim:" "`
	}
	RegisterTaskFatal("git_gc_repos", &RepoHealthCheckConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@every 72h",
		},
		Timeout: time.Duration(setting.Git.Timeout.GC) * time.Second,
		Args:    setting.Git.GCArgs,
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		rhcConfig := config.(*RepoHealthCheckConfig)
		// the git args are set by config, they can be safe to be trusted
		return repo_service.GitGcRepos(ctx, rhcConfig.Timeout, gitcmd.ToTrustedCmdArgs(rhcConfig.Args))
	})
}

func registerRewriteAllPublicKeys() {
	RegisterTaskFatal("resync_all_sshkeys", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@every 72h",
	}, func(ctx context.Context, _ *user_model.User, _ Config) error {
		return asymkey_service.RewriteAllPublicKeys(ctx)
	})
}

func registerRewriteAllPrincipalKeys() {
	RegisterTaskFatal("resync_all_sshprincipals", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@every 72h",
	}, func(ctx context.Context, _ *user_model.User, _ Config) error {
		return asymkey_service.RewriteAllPrincipalKeys(ctx)
	})
}

func registerRepositoryUpdateHook() {
	RegisterTaskFatal("resync_all_hooks", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@every 72h",
	}, func(ctx context.Context, _ *user_model.User, _ Config) error {
		return repo_service.SyncRepositoryHooks(ctx)
	})
}

func registerReinitMissingRepositories() {
	RegisterTaskFatal("reinit_missing_repos", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@every 72h",
	}, func(ctx context.Context, _ *user_model.User, _ Config) error {
		return repo_service.ReinitMissingRepositories(ctx)
	})
}

func registerDeleteMissingRepositories() {
	RegisterTaskFatal("delete_missing_repos", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@every 72h",
	}, func(ctx context.Context, user *user_model.User, _ Config) error {
		return repo_service.DeleteMissingRepositories(ctx, user)
	})
}

func registerRemoveRandomAvatars() {
	RegisterTaskFatal("delete_generated_repository_avatars", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@every 72h",
	}, func(ctx context.Context, _ *user_model.User, _ Config) error {
		return repo_service.RemoveRandomAvatars(ctx)
	})
}

func registerDeleteOldActions() {
	RegisterTaskFatal("delete_old_actions", &OlderThanConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@every 168h",
		},
		OlderThan: 365 * 24 * time.Hour,
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		olderThanConfig := config.(*OlderThanConfig)
		return activities_model.DeleteOldActions(ctx, olderThanConfig.OlderThan)
	})
}

func registerUpdateKmupChecker() {
	type UpdateCheckerConfig struct {
		BaseConfig
		HTTPEndpoint string
	}
	RegisterTaskFatal("update_checker", &UpdateCheckerConfig{
		BaseConfig: BaseConfig{
			Enabled:    true,
			RunAtStart: false,
			Schedule:   "@every 168h",
		},
		HTTPEndpoint: "https://dl.kmup.com/kmup/version.json",
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		updateCheckerConfig := config.(*UpdateCheckerConfig)
		return updatechecker.KmupUpdateChecker(updateCheckerConfig.HTTPEndpoint)
	})
}

func registerDeleteOldSystemNotices() {
	RegisterTaskFatal("delete_old_system_notices", &OlderThanConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@every 168h",
		},
		OlderThan: 365 * 24 * time.Hour,
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		olderThanConfig := config.(*OlderThanConfig)
		return system.DeleteOldSystemNotices(ctx, olderThanConfig.OlderThan)
	})
}

type GCLFSConfig struct {
	BaseConfig
	OlderThan                time.Duration
	LastUpdatedMoreThanAgo   time.Duration
	NumberToCheckPerRepo     int64
	ProportionToCheckPerRepo float64
}

func registerGCLFS() {
	if !setting.LFS.StartServer {
		return
	}

	RegisterTaskFatal("gc_lfs", &GCLFSConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@every 24h",
		},
		// Only attempt to garbage collect lfs meta objects older than a week as the order of git lfs upload
		// and git object upload is not necessarily guaranteed. It's possible to imagine a situation whereby
		// an LFS object is uploaded but the git branch is not uploaded immediately, or there are some rapid
		// changes in new branches that might lead to lfs objects becoming temporarily unassociated with git
		// objects.
		//
		// It is likely that a week is potentially excessive but it should definitely be enough that any
		// unassociated LFS object is genuinely unassociated.
		OlderThan: 24 * time.Hour * 7,

		// Only GC things that haven't been looked at in the past 3 days
		LastUpdatedMoreThanAgo:   24 * time.Hour * 3,
		NumberToCheckPerRepo:     100,
		ProportionToCheckPerRepo: 0.6,
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		gcLFSConfig := config.(*GCLFSConfig)
		return repo_service.GarbageCollectLFSMetaObjects(ctx, repo_service.GarbageCollectLFSMetaObjectsOptions{
			AutoFix:                 true,
			OlderThan:               time.Now().Add(-gcLFSConfig.OlderThan),
			UpdatedLessRecentlyThan: time.Now().Add(-gcLFSConfig.LastUpdatedMoreThanAgo),
		})
	})
}

func registerRebuildIssueIndexer() {
	RegisterTaskFatal("rebuild_issue_indexer", &BaseConfig{
		Enabled:    false,
		RunAtStart: false,
		Schedule:   "@annually",
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		return issue_indexer.PopulateIssueIndexer(ctx)
	})
}

func initExtendedTasks() {
	registerDeleteInactiveUsers()
	registerDeleteRepositoryArchives()
	registerGarbageCollectRepositories()
	registerRewriteAllPublicKeys()
	registerRewriteAllPrincipalKeys()
	registerRepositoryUpdateHook()
	registerReinitMissingRepositories()
	registerDeleteMissingRepositories()
	registerRemoveRandomAvatars()
	registerDeleteOldActions()
	registerUpdateKmupChecker()
	registerDeleteOldSystemNotices()
	registerGCLFS()
	registerRebuildIssueIndexer()
}
