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

package admin

import (
	"net/http"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/cron"
)

// ListCronTasks api for getting cron tasks
func ListCronTasks(ctx *context.APIContext) {
	// swagger:operation GET /admin/cron admin adminCronList
	// ---
	// summary: List cron tasks
	// produces:
	// - application/json
	// parameters:
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/CronList"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	tasks := cron.ListTasks()
	count := len(tasks)

	listOpts := utils.GetListOptions(ctx)
	tasks = util.PaginateSlice(tasks, listOpts.Page, listOpts.PageSize).(cron.TaskTable)

	res := make([]structs.Cron, len(tasks))
	for i, task := range tasks {
		res[i] = structs.Cron{
			Name:      task.Name,
			Schedule:  task.Spec,
			Next:      task.Next,
			Prev:      task.Prev,
			ExecTimes: task.ExecTimes,
		}
	}

	ctx.SetTotalCountHeader(int64(count))
	ctx.JSON(http.StatusOK, res)
}

// PostCronTask api for getting cron tasks
func PostCronTask(ctx *context.APIContext) {
	// swagger:operation POST /admin/cron/{task} admin adminCronRun
	// ---
	// summary: Run cron task
	// produces:
	// - application/json
	// parameters:
	// - name: task
	//   in: path
	//   description: task to run
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "404":
	//     "$ref": "#/responses/notFound"
	task := cron.GetTask(ctx.PathParam("task"))
	if task == nil {
		ctx.APIErrorNotFound()
		return
	}
	task.Run()
	log.Trace("Cron Task %s started by admin(%s)", task.Name, ctx.Doer.Name)

	ctx.Status(http.StatusNoContent)
}
