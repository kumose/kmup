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
	"errors"
	"fmt"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/models/db"
	secret_model "github.com/kumose/kmup/models/secret"
	notify_service "github.com/kumose/kmup/services/notify"

	runnerv1 "github.com/kumose/actions-proto-go/runner/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func PickTask(ctx context.Context, runner *actions_model.ActionRunner) (*runnerv1.Task, bool, error) {
	var (
		task       *runnerv1.Task
		job        *actions_model.ActionRunJob
		actionTask *actions_model.ActionTask
	)

	if runner.Ephemeral {
		var task actions_model.ActionTask
		has, err := db.GetEngine(ctx).Where("runner_id = ?", runner.ID).Get(&task)
		// Let the runner retry the request, do not allow to proceed
		if err != nil {
			return nil, false, err
		}
		if has {
			if task.Status == actions_model.StatusWaiting || task.Status == actions_model.StatusRunning || task.Status == actions_model.StatusBlocked {
				return nil, false, nil
			}
			// task has been finished, remove it
			_, err = db.DeleteByID[actions_model.ActionRunner](ctx, runner.ID)
			if err != nil {
				return nil, false, err
			}
			return nil, false, errors.New("runner has been removed")
		}
	}

	if err := db.WithTx(ctx, func(ctx context.Context) error {
		t, ok, err := actions_model.CreateTaskForRunner(ctx, runner)
		if err != nil {
			return fmt.Errorf("CreateTaskForRunner: %w", err)
		}
		if !ok {
			return nil
		}

		if err := t.LoadAttributes(ctx); err != nil {
			return fmt.Errorf("task LoadAttributes: %w", err)
		}
		job = t.Job
		actionTask = t

		secrets, err := secret_model.GetSecretsOfTask(ctx, t)
		if err != nil {
			return fmt.Errorf("GetSecretsOfTask: %w", err)
		}

		vars, err := actions_model.GetVariablesOfRun(ctx, t.Job.Run)
		if err != nil {
			return fmt.Errorf("GetVariablesOfRun: %w", err)
		}

		needs, err := findTaskNeeds(ctx, job)
		if err != nil {
			return fmt.Errorf("findTaskNeeds: %w", err)
		}

		taskContext, err := generateTaskContext(t)
		if err != nil {
			return fmt.Errorf("generateTaskContext: %w", err)
		}

		task = &runnerv1.Task{
			Id:              t.ID,
			WorkflowPayload: t.Job.WorkflowPayload,
			Context:         taskContext,
			Secrets:         secrets,
			Vars:            vars,
			Needs:           needs,
		}

		return nil
	}); err != nil {
		return nil, false, err
	}

	if task == nil {
		return nil, false, nil
	}

	CreateCommitStatusForRunJobs(ctx, job.Run, job)
	notify_service.WorkflowJobStatusUpdate(ctx, job.Run.Repo, job.Run.TriggerUser, job, actionTask)

	return task, true, nil
}

func generateTaskContext(t *actions_model.ActionTask) (*structpb.Struct, error) {
	kmupRuntimeToken, err := CreateAuthorizationToken(t.ID, t.Job.RunID, t.JobID)
	if err != nil {
		return nil, err
	}

	gitCtx := GenerateKmupContext(t.Job.Run, t.Job)
	gitCtx["token"] = t.Token
	gitCtx["kmup_runtime_token"] = kmupRuntimeToken

	return structpb.NewStruct(gitCtx)
}

func findTaskNeeds(ctx context.Context, taskJob *actions_model.ActionRunJob) (map[string]*runnerv1.TaskNeed, error) {
	taskNeeds, err := FindTaskNeeds(ctx, taskJob)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*runnerv1.TaskNeed, len(taskNeeds))
	for jobID, taskNeed := range taskNeeds {
		ret[jobID] = &runnerv1.TaskNeed{
			Outputs: taskNeed.Outputs,
			Result:  runnerv1.Result(taskNeed.Result),
		}
	}
	return ret, nil
}
