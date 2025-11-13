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

package release

import (
	"context"
	"errors"
	"fmt"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/queue"
	repo_module "github.com/kumose/kmup/modules/repository"

	"xorm.io/builder"
)

type TagSyncOptions struct {
	RepoID int64
}

// tagSyncQueue represents a queue to handle tag sync jobs.
var tagSyncQueue *queue.WorkerPoolQueue[*TagSyncOptions]

func handlerTagSync(items ...*TagSyncOptions) []*TagSyncOptions {
	for _, opts := range items {
		err := repo_module.SyncRepoTags(graceful.GetManager().ShutdownContext(), opts.RepoID)
		if err != nil {
			log.Error("syncRepoTags [%d] failed: %v", opts.RepoID, err)
		}
	}
	return nil
}

func addRepoToTagSyncQueue(repoID int64) error {
	return tagSyncQueue.Push(&TagSyncOptions{
		RepoID: repoID,
	})
}

func initTagSyncQueue(ctx context.Context) error {
	tagSyncQueue = queue.CreateUniqueQueue(ctx, "tag_sync", handlerTagSync)
	if tagSyncQueue == nil {
		return errors.New("unable to create tag_sync queue")
	}
	go graceful.GetManager().RunWithCancel(tagSyncQueue)

	return nil
}

func AddAllRepoTagsToSyncQueue(ctx context.Context) error {
	if err := db.Iterate(ctx, builder.Eq{"is_empty": false}, func(ctx context.Context, repo *repo_model.Repository) error {
		return addRepoToTagSyncQueue(repo.ID)
	}); err != nil {
		return fmt.Errorf("run sync all tags failed: %v", err)
	}
	return nil
}
