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

package doctor

import (
	"context"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/storage"
	repo_service "github.com/kumose/kmup/services/repository"

	"xorm.io/builder"
)

func handleDeleteOrphanedRepos(ctx context.Context, logger log.Logger, autofix bool) error {
	test := &consistencyCheck{
		Name:         "Repos with no existing owner",
		Counter:      countOrphanedRepos,
		Fixer:        deleteOrphanedRepos,
		FixedMessage: "Deleted all content related to orphaned repos",
	}
	return test.Run(ctx, logger, autofix)
}

// countOrphanedRepos count repository where user of owner_id do not exist
func countOrphanedRepos(ctx context.Context) (int64, error) {
	return db.CountOrphanedObjects(ctx, "repository", "user", "repository.owner_id=`user`.id")
}

// deleteOrphanedRepos delete repository where user of owner_id do not exist
func deleteOrphanedRepos(ctx context.Context) (int64, error) {
	if err := storage.Init(); err != nil {
		return 0, err
	}

	batchSize := db.MaxBatchInsertSize("repository")
	e := db.GetEngine(ctx)
	var deleted int64

	for {
		select {
		case <-ctx.Done():
			return deleted, ctx.Err()
		default:
			var ids []int64
			if err := e.Table("`repository`").
				Join("LEFT", "`user`", "repository.owner_id=`user`.id").
				Where(builder.IsNull{"`user`.id"}).
				Select("`repository`.id").Limit(batchSize).Find(&ids); err != nil {
				return deleted, err
			}

			// if we don't get ids we have deleted them all
			if len(ids) == 0 {
				return deleted, nil
			}

			for _, id := range ids {
				if err := repo_service.DeleteRepositoryDirectly(ctx, id, true); err != nil {
					return deleted, err
				}
				deleted++
			}
		}
	}
}

func init() {
	Register(&Check{
		Title:     "Deleted all content related to orphaned repos",
		Name:      "delete-orphaned-repos",
		IsDefault: false,
		Run:       handleDeleteOrphanedRepos,
		Priority:  4,
	})
}
