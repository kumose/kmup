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

package common

import (
	"fmt"
	"strings"

	actions_model "github.com/kumose/kmup/models/actions"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/actions"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
)

func DownloadActionsRunJobLogsWithIndex(ctx *context.Base, ctxRepo *repo_model.Repository, runID, jobIndex int64) error {
	runJobs, err := actions_model.GetRunJobsByRunID(ctx, runID)
	if err != nil {
		return fmt.Errorf("GetRunJobsByRunID: %w", err)
	}
	if err = runJobs.LoadRepos(ctx); err != nil {
		return fmt.Errorf("LoadRepos: %w", err)
	}
	if jobIndex < 0 || jobIndex >= int64(len(runJobs)) {
		return util.NewNotExistErrorf("job index is out of range: %d", jobIndex)
	}
	return DownloadActionsRunJobLogs(ctx, ctxRepo, runJobs[jobIndex])
}

func DownloadActionsRunJobLogs(ctx *context.Base, ctxRepo *repo_model.Repository, curJob *actions_model.ActionRunJob) error {
	if curJob.Repo.ID != ctxRepo.ID {
		return util.NewNotExistErrorf("job not found")
	}

	if curJob.TaskID == 0 {
		return util.NewNotExistErrorf("job not started")
	}

	if err := curJob.LoadRun(ctx); err != nil {
		return fmt.Errorf("LoadRun: %w", err)
	}

	task, err := actions_model.GetTaskByID(ctx, curJob.TaskID)
	if err != nil {
		return fmt.Errorf("GetTaskByID: %w", err)
	}

	if task.LogExpired {
		return util.NewNotExistErrorf("logs have been cleaned up")
	}

	reader, err := actions.OpenLogs(ctx, task.LogInStorage, task.LogFilename)
	if err != nil {
		return fmt.Errorf("OpenLogs: %w", err)
	}
	defer reader.Close()

	workflowName := curJob.Run.WorkflowID
	if p := strings.Index(workflowName, "."); p > 0 {
		workflowName = workflowName[0:p]
	}
	ctx.ServeContent(reader, &context.ServeHeaderOptions{
		Filename:           fmt.Sprintf("%v-%v-%v.log", workflowName, curJob.Name, task.ID),
		ContentLength:      &task.LogSize,
		ContentType:        "text/plain",
		ContentTypeCharset: "utf-8",
		Disposition:        "attachment",
	})
	return nil
}
