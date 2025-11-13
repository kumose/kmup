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

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/webhook"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/log"

	"xorm.io/builder"
)

// SyncRepositoryHooks rewrites all repositories' pre-receive, update and post-receive hooks
// to make sure the binary and custom conf path are up-to-date.
func SyncRepositoryHooks(ctx context.Context) error {
	log.Trace("Doing: SyncRepositoryHooks")

	if err := db.Iterate(
		ctx,
		builder.Gt{"id": 0},
		func(ctx context.Context, repo *repo_model.Repository) error {
			select {
			case <-ctx.Done():
				return db.ErrCancelledf("before sync repository hooks for %s", repo.FullName())
			default:
			}

			if err := gitrepo.CreateDelegateHooks(ctx, repo); err != nil {
				return fmt.Errorf("CreateDelegateHooks: %w", err)
			}

			if HasWiki(ctx, repo) {
				if err := gitrepo.CreateDelegateHooks(ctx, repo.WikiStorageRepo()); err != nil {
					return fmt.Errorf("CreateDelegateHooks: %w", err)
				}
			}
			return nil
		},
	); err != nil {
		return err
	}

	log.Trace("Finished: SyncRepositoryHooks")
	return nil
}

// GenerateGitHooks generates git hooks from a template repository
func GenerateGitHooks(ctx context.Context, templateRepo, generateRepo *repo_model.Repository) error {
	generateGitRepo, err := gitrepo.OpenRepository(ctx, generateRepo)
	if err != nil {
		return err
	}
	defer generateGitRepo.Close()

	templateGitRepo, err := gitrepo.OpenRepository(ctx, templateRepo)
	if err != nil {
		return err
	}
	defer templateGitRepo.Close()

	templateHooks, err := templateGitRepo.Hooks()
	if err != nil {
		return err
	}

	for _, templateHook := range templateHooks {
		generateHook, err := generateGitRepo.GetHook(templateHook.Name())
		if err != nil {
			return err
		}

		generateHook.Content = templateHook.Content
		if err := generateHook.Update(); err != nil {
			return err
		}
	}
	return nil
}

// GenerateWebhooks generates webhooks from a template repository
func GenerateWebhooks(ctx context.Context, templateRepo, generateRepo *repo_model.Repository) error {
	templateWebhooks, err := db.Find[webhook.Webhook](ctx, webhook.ListWebhookOptions{RepoID: templateRepo.ID})
	if err != nil {
		return err
	}

	ws := make([]*webhook.Webhook, 0, len(templateWebhooks))
	for _, templateWebhook := range templateWebhooks {
		ws = append(ws, &webhook.Webhook{
			RepoID:      generateRepo.ID,
			URL:         templateWebhook.URL,
			HTTPMethod:  templateWebhook.HTTPMethod,
			ContentType: templateWebhook.ContentType,
			Secret:      templateWebhook.Secret,
			HookEvent:   templateWebhook.HookEvent,
			IsActive:    templateWebhook.IsActive,
			Type:        templateWebhook.Type,
			OwnerID:     templateWebhook.OwnerID,
			Events:      templateWebhook.Events,
			Meta:        templateWebhook.Meta,
		})
	}
	return webhook.CreateWebhooks(ctx, ws)
}
