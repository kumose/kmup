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

package repo

import (
	"errors"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/routers/common"
	"github.com/kumose/kmup/services/context"
)

func DownloadActionsRunJobLogs(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/actions/jobs/{job_id}/logs repository downloadActionsRunJobLogs
	// ---
	// summary: Downloads the job logs for a workflow run
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repository
	//   type: string
	//   required: true
	// - name: job_id
	//   in: path
	//   description: id of the job
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     description: output blob content
	//   "400":
	//     "$ref": "#/responses/error"
	//   "404":
	//     "$ref": "#/responses/notFound"

	jobID := ctx.PathParamInt64("job_id")
	curJob, err := actions_model.GetRunJobByID(ctx, jobID)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	if err = curJob.LoadRepo(ctx); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	err = common.DownloadActionsRunJobLogs(ctx.Base, ctx.Repo.Repository, curJob)
	if err != nil {
		if errors.Is(err, util.ErrNotExist) {
			ctx.APIErrorNotFound(err)
		} else {
			ctx.APIErrorInternal(err)
		}
	}
}
