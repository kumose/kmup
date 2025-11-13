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
	"fmt"
	"time"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"
	webhook_module "github.com/kumose/kmup/modules/webhook"
)

// ActionSchedule represents a schedule of a workflow file
type ActionSchedule struct {
	ID            int64
	Title         string
	Specs         []string
	RepoID        int64                  `xorm:"index"`
	Repo          *repo_model.Repository `xorm:"-"`
	OwnerID       int64                  `xorm:"index"`
	WorkflowID    string
	TriggerUserID int64
	TriggerUser   *user_model.User `xorm:"-"`
	Ref           string
	CommitSHA     string
	Event         webhook_module.HookEventType
	EventPayload  string `xorm:"LONGTEXT"`
	Content       []byte
	Created       timeutil.TimeStamp `xorm:"created"`
	Updated       timeutil.TimeStamp `xorm:"updated"`
}

func init() {
	db.RegisterModel(new(ActionSchedule))
}

// GetSchedulesMapByIDs returns the schedules by given id slice.
func GetSchedulesMapByIDs(ctx context.Context, ids []int64) (map[int64]*ActionSchedule, error) {
	schedules := make(map[int64]*ActionSchedule, len(ids))
	if len(ids) == 0 {
		return schedules, nil
	}
	return schedules, db.GetEngine(ctx).In("id", ids).Find(&schedules)
}

// CreateScheduleTask creates new schedule task.
func CreateScheduleTask(ctx context.Context, rows []*ActionSchedule) error {
	// Return early if there are no rows to insert
	if len(rows) == 0 {
		return nil
	}

	return db.WithTx(ctx, func(ctx context.Context) error {
		// Loop through each schedule row
		for _, row := range rows {
			row.Title = util.EllipsisDisplayString(row.Title, 255)
			// Create new schedule row
			if err := db.Insert(ctx, row); err != nil {
				return err
			}

			// Loop through each schedule spec and create a new spec row
			now := time.Now()

			for _, spec := range row.Specs {
				specRow := &ActionScheduleSpec{
					RepoID:     row.RepoID,
					ScheduleID: row.ID,
					Spec:       spec,
				}
				// Parse the spec and check for errors
				schedule, err := specRow.Parse()
				if err != nil {
					continue // skip to the next spec if there's an error
				}

				specRow.Next = timeutil.TimeStamp(schedule.Next(now).Unix())

				// Insert the new schedule spec row
				if err = db.Insert(ctx, specRow); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func DeleteScheduleTaskByRepo(ctx context.Context, id int64) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		if _, err := db.GetEngine(ctx).Delete(&ActionSchedule{RepoID: id}); err != nil {
			return err
		}

		if _, err := db.GetEngine(ctx).Delete(&ActionScheduleSpec{RepoID: id}); err != nil {
			return err
		}

		return nil
	})
}

func CleanRepoScheduleTasks(ctx context.Context, repo *repo_model.Repository) ([]*ActionRunJob, error) {
	// If actions disabled when there is schedule task, this will remove the outdated schedule tasks
	// There is no other place we can do this because the app.ini will be changed manually
	if err := DeleteScheduleTaskByRepo(ctx, repo.ID); err != nil {
		return nil, fmt.Errorf("DeleteCronTaskByRepo: %v", err)
	}
	// cancel running cron jobs of this repository and delete old schedules
	jobs, err := CancelPreviousJobs(
		ctx,
		repo.ID,
		repo.DefaultBranch,
		"",
		webhook_module.HookEventSchedule,
	)
	if err != nil {
		return jobs, fmt.Errorf("CancelPreviousJobs: %v", err)
	}
	return jobs, nil
}
