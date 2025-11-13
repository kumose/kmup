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
	"runtime"

	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/context"
)

func monitorTraceCommon(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.monitor")
	ctx.Data["PageIsAdminMonitorTrace"] = true
	// Hide the performance trace tab in production, because it shows a lot of SQLs and is not that useful for end users.
	// To avoid confusing end users, do not let them know this tab. End users should "download diagnosis report" instead.
	ctx.Data["ShowAdminPerformanceTraceTab"] = !setting.IsProd
}

// Stacktrace show admin monitor goroutines page
func Stacktrace(ctx *context.Context) {
	monitorTraceCommon(ctx)

	ctx.Data["GoroutineCount"] = runtime.NumGoroutine()

	show := ctx.FormString("show")
	ctx.Data["ShowGoroutineList"] = show
	// by default, do not do anything which might cause server errors, to avoid unnecessary 500 pages.
	// this page is the entrance of the chance to collect diagnosis report.
	if show != "" {
		showNoSystem := show == "process"
		processStacks, processCount, _, err := process.GetManager().ProcessStacktraces(false, showNoSystem)
		if err != nil {
			ctx.ServerError("GoroutineStacktrace", err)
			return
		}

		ctx.Data["ProcessStacks"] = processStacks
		ctx.Data["ProcessCount"] = processCount
	}

	ctx.HTML(http.StatusOK, tplStacktrace)
}

// StacktraceCancel cancels a process
func StacktraceCancel(ctx *context.Context) {
	pid := ctx.PathParam("pid")
	process.GetManager().Cancel(process.IDType(pid))
	ctx.JSONRedirect(setting.AppSubURL + "/-/admin/monitor/stacktrace")
}
