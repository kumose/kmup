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

package actions

import (
	"context"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/container"
)

type RunnerList []*ActionRunner

// GetUserIDs returns a slice of user's id
func (runners RunnerList) GetUserIDs() []int64 {
	return container.FilterSlice(runners, func(runner *ActionRunner) (int64, bool) {
		return runner.OwnerID, runner.OwnerID != 0
	})
}

func (runners RunnerList) LoadOwners(ctx context.Context) error {
	userIDs := runners.GetUserIDs()
	users := make(map[int64]*user_model.User, len(userIDs))
	if err := db.GetEngine(ctx).In("id", userIDs).Find(&users); err != nil {
		return err
	}
	for _, runner := range runners {
		if runner.OwnerID > 0 && runner.Owner == nil {
			runner.Owner = users[runner.OwnerID]
		}
	}
	return nil
}

func (runners RunnerList) getRepoIDs() []int64 {
	return container.FilterSlice(runners, func(runner *ActionRunner) (int64, bool) {
		return runner.RepoID, runner.RepoID > 0
	})
}

func (runners RunnerList) LoadRepos(ctx context.Context) error {
	repoIDs := runners.getRepoIDs()
	repos := make(map[int64]*repo_model.Repository, len(repoIDs))
	if err := db.GetEngine(ctx).In("id", repoIDs).Find(&repos); err != nil {
		return err
	}

	for _, runner := range runners {
		if runner.RepoID > 0 && runner.Repo == nil {
			runner.Repo = repos[runner.RepoID]
		}
	}
	return nil
}

func (runners RunnerList) LoadAttributes(ctx context.Context) error {
	if err := runners.LoadOwners(ctx); err != nil {
		return err
	}

	return runners.LoadRepos(ctx)
}
