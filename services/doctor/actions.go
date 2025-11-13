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
	"fmt"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	unit_model "github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/timeutil"
	repo_service "github.com/kumose/kmup/services/repository"

	"xorm.io/builder"
)

func disableMirrorActionsUnit(ctx context.Context, logger log.Logger, autofix bool) error {
	var reposToFix []*repo_model.Repository

	for page := 1; ; page++ {
		repos, _, err := repo_model.SearchRepository(ctx, repo_model.SearchRepoOptions{
			ListOptions: db.ListOptions{
				PageSize: repo_model.RepositoryListDefaultPageSize,
				Page:     page,
			},
			Mirror: optional.Some(true),
		})
		if err != nil {
			return fmt.Errorf("SearchRepository: %w", err)
		}
		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			if repo.UnitEnabled(ctx, unit_model.TypeActions) {
				reposToFix = append(reposToFix, repo)
			}
		}
	}

	if len(reposToFix) == 0 {
		logger.Info("Found no mirror with actions unit enabled")
	} else {
		logger.Warn("Found %d mirrors with actions unit enabled", len(reposToFix))
	}
	if !autofix || len(reposToFix) == 0 {
		return nil
	}

	for _, repo := range reposToFix {
		if err := repo_service.UpdateRepositoryUnits(ctx, repo, nil, []unit_model.Type{unit_model.TypeActions}); err != nil {
			return err
		}
	}
	logger.Info("Fixed %d mirrors with actions unit enabled", len(reposToFix))

	return nil
}

func fixUnfinishedRunStatus(ctx context.Context, logger log.Logger, autofix bool) error {
	total := 0
	inconsistent := 0
	fixed := 0

	cond := builder.In("status", []actions_model.Status{
		actions_model.StatusWaiting,
		actions_model.StatusRunning,
		actions_model.StatusBlocked,
	}).And(builder.Lt{"updated": timeutil.TimeStampNow().AddDuration(-setting.Actions.ZombieTaskTimeout)})

	err := db.Iterate(
		ctx,
		cond,
		func(ctx context.Context, run *actions_model.ActionRun) error {
			total++

			jobs, err := actions_model.GetRunJobsByRunID(ctx, run.ID)
			if err != nil {
				return fmt.Errorf("GetRunJobsByRunID: %w", err)
			}
			expected := actions_model.AggregateJobStatus(jobs)
			if expected == run.Status {
				return nil
			}

			inconsistent++
			logger.Warn("Run %d (repo_id=%d, index=%d) has status %s, expected %s", run.ID, run.RepoID, run.Index, run.Status, expected)

			if !autofix {
				return nil
			}

			run.Started, run.Stopped = getRunTimestampsFromJobs(run, expected, jobs)
			run.Status = expected

			if err := actions_model.UpdateRun(ctx, run, "status", "started", "stopped"); err != nil {
				return fmt.Errorf("UpdateRun: %w", err)
			}
			fixed++

			return nil
		},
	)
	if err != nil {
		logger.Critical("Unable to iterate unfinished runs: %v", err)
		return err
	}

	if inconsistent == 0 {
		logger.Info("Checked %d unfinished runs; all statuses are consistent.", total)
		return nil
	}

	if autofix {
		logger.Info("Checked %d unfinished runs; fixed %d of %d runs.", total, fixed, inconsistent)
	} else {
		logger.Warn("Checked %d unfinished runs; found %d runs need to be fixed", total, inconsistent)
	}

	return nil
}

func getRunTimestampsFromJobs(run *actions_model.ActionRun, newStatus actions_model.Status, jobs actions_model.ActionJobList) (started, stopped timeutil.TimeStamp) {
	started = run.Started
	if (newStatus.IsRunning() || newStatus.IsDone()) && started.IsZero() {
		var earliest timeutil.TimeStamp
		for _, job := range jobs {
			if job.Started > 0 && (earliest.IsZero() || job.Started < earliest) {
				earliest = job.Started
			}
		}
		started = earliest
	}

	stopped = run.Stopped
	if newStatus.IsDone() && stopped.IsZero() {
		var latest timeutil.TimeStamp
		for _, job := range jobs {
			if job.Stopped > latest {
				latest = job.Stopped
			}
		}
		stopped = latest
	}

	return started, stopped
}

func init() {
	Register(&Check{
		Title:     "Disable the actions unit for all mirrors",
		Name:      "disable-mirror-actions-unit",
		IsDefault: false,
		Run:       disableMirrorActionsUnit,
		Priority:  9,
	})
	Register(&Check{
		Title:     "Fix inconsistent status for unfinished actions runs",
		Name:      "fix-actions-unfinished-run-status",
		IsDefault: false,
		Run:       fixUnfinishedRunStatus,
		Priority:  9,
	})
}
