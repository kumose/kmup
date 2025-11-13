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

	"xorm.io/builder"
)

type ScheduleList []*ActionSchedule

// GetUserIDs returns a slice of user's id
func (schedules ScheduleList) GetUserIDs() []int64 {
	return container.FilterSlice(schedules, func(schedule *ActionSchedule) (int64, bool) {
		return schedule.TriggerUserID, true
	})
}

func (schedules ScheduleList) GetRepoIDs() []int64 {
	return container.FilterSlice(schedules, func(schedule *ActionSchedule) (int64, bool) {
		return schedule.RepoID, true
	})
}

func (schedules ScheduleList) LoadTriggerUser(ctx context.Context) error {
	userIDs := schedules.GetUserIDs()
	users := make(map[int64]*user_model.User, len(userIDs))
	if err := db.GetEngine(ctx).In("id", userIDs).Find(&users); err != nil {
		return err
	}
	for _, schedule := range schedules {
		if schedule.TriggerUserID == user_model.ActionsUserID {
			schedule.TriggerUser = user_model.NewActionsUser()
		} else {
			schedule.TriggerUser = users[schedule.TriggerUserID]
			if schedule.TriggerUser == nil {
				schedule.TriggerUser = user_model.NewGhostUser()
			}
		}
	}
	return nil
}

func (schedules ScheduleList) LoadRepos(ctx context.Context) error {
	repoIDs := schedules.GetRepoIDs()
	repos, err := repo_model.GetRepositoriesMapByIDs(ctx, repoIDs)
	if err != nil {
		return err
	}
	for _, schedule := range schedules {
		schedule.Repo = repos[schedule.RepoID]
	}
	return nil
}

type FindScheduleOptions struct {
	db.ListOptions
	RepoID  int64
	OwnerID int64
}

func (opts FindScheduleOptions) ToConds() builder.Cond {
	cond := builder.NewCond()
	if opts.RepoID > 0 {
		cond = cond.And(builder.Eq{"repo_id": opts.RepoID})
	}
	if opts.OwnerID > 0 {
		cond = cond.And(builder.Eq{"owner_id": opts.OwnerID})
	}

	return cond
}

func (opts FindScheduleOptions) ToOrders() string {
	return "`id` DESC"
}
