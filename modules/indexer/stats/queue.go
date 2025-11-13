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

package stats

import (
	"errors"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/queue"
	"github.com/kumose/kmup/modules/setting"
)

// statsQueue represents a queue to handle repository stats updates
var statsQueue *queue.WorkerPoolQueue[int64]

// handle passed PR IDs and test the PRs
func handler(items ...int64) []int64 {
	for _, opts := range items {
		if err := indexer.Index(opts); err != nil {
			if !setting.IsInTesting {
				log.Error("stats queue indexer.Index(%d) failed: %v", opts, err)
			}
		}
	}
	return nil
}

func initStatsQueue() error {
	statsQueue = queue.CreateUniqueQueue(graceful.GetManager().ShutdownContext(), "repo_stats_update", handler)
	if statsQueue == nil {
		return errors.New("unable to create repo_stats_update queue")
	}
	go graceful.GetManager().RunWithCancel(statsQueue)
	return nil
}

// UpdateRepoIndexer update a repository's entries in the indexer
func UpdateRepoIndexer(repo *repo_model.Repository) error {
	if err := statsQueue.Push(repo.ID); err != nil {
		if err != queue.ErrAlreadyInQueue {
			return err
		}
		log.Debug("Repo ID: %d already queued", repo.ID)
	}
	return nil
}
