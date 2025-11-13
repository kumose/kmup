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

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	actions_service "github.com/kumose/kmup/services/actions"
)

func initActionsTasks() {
	if !setting.Actions.Enabled {
		return
	}
	registerStopZombieTasks()
	registerStopEndlessTasks()
	registerCancelAbandonedJobs()
	registerScheduleTasks()
	registerActionsCleanup()
}

func registerStopZombieTasks() {
	RegisterTaskFatal("stop_zombie_tasks", &BaseConfig{
		Enabled:    true,
		RunAtStart: true,
		Schedule:   "@every 5m",
	}, func(ctx context.Context, _ *user_model.User, cfg Config) error {
		return actions_service.StopZombieTasks(ctx)
	})
}

func registerStopEndlessTasks() {
	RegisterTaskFatal("stop_endless_tasks", &BaseConfig{
		Enabled:    true,
		RunAtStart: true,
		Schedule:   "@every 30m",
	}, func(ctx context.Context, _ *user_model.User, cfg Config) error {
		return actions_service.StopEndlessTasks(ctx)
	})
}

func registerCancelAbandonedJobs() {
	RegisterTaskFatal("cancel_abandoned_jobs", &BaseConfig{
		Enabled:    true,
		RunAtStart: true,
		Schedule:   "@every 6h",
	}, func(ctx context.Context, _ *user_model.User, cfg Config) error {
		return actions_service.CancelAbandonedJobs(ctx)
	})
}

// registerScheduleTasks registers a scheduled task that runs every minute to start any due schedule tasks.
func registerScheduleTasks() {
	// Register the task with a unique name, enabled status, and schedule for every minute.
	RegisterTaskFatal("start_schedule_tasks", &BaseConfig{
		Enabled:    true,
		RunAtStart: false,
		Schedule:   "@every 1m",
	}, func(ctx context.Context, _ *user_model.User, cfg Config) error {
		// Call the function to start schedule tasks and pass the context.
		return actions_service.StartScheduleTasks(ctx)
	})
}

func registerActionsCleanup() {
	RegisterTaskFatal("cleanup_actions", &BaseConfig{
		Enabled:    true,
		RunAtStart: false,
		Schedule:   "@midnight",
	}, func(ctx context.Context, _ *user_model.User, _ Config) error {
		return actions_service.Cleanup(ctx)
	})
}
