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

package automergequeue

import (
	"context"
	"errors"
	"fmt"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/queue"
)

var AutoMergeQueue *queue.WorkerPoolQueue[string]

var AddToQueue = func(pr *issues_model.PullRequest, sha string) {
	log.Trace("Adding pullID: %d to the pull requests patch checking queue with sha %s", pr.ID, sha)
	if err := AutoMergeQueue.Push(fmt.Sprintf("%d_%s", pr.ID, sha)); err != nil && !errors.Is(err, queue.ErrAlreadyInQueue) {
		log.Error("Error adding pullID: %d to the pull requests patch checking queue %v", pr.ID, err)
	}
}

// StartPRCheckAndAutoMerge start an automerge check and auto merge task for a pull request
func StartPRCheckAndAutoMerge(ctx context.Context, pull *issues_model.PullRequest) {
	if pull == nil || pull.HasMerged || !pull.CanAutoMerge() {
		return
	}

	if err := pull.LoadBaseRepo(ctx); err != nil {
		log.Error("LoadBaseRepo: %v", err)
		return
	}

	gitRepo, err := gitrepo.OpenRepository(ctx, pull.BaseRepo)
	if err != nil {
		log.Error("OpenRepository: %v", err)
		return
	}
	defer gitRepo.Close()
	commitID, err := gitRepo.GetRefCommitID(pull.GetGitHeadRefName())
	if err != nil {
		log.Error("GetRefCommitID: %v", err)
		return
	}

	AddToQueue(pull, commitID)
}
