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
	"strconv"

	"github.com/kumose/kmup/modules/queue"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/context"
)

func Queues(ctx *context.Context) {
	if !setting.IsProd {
		initTestQueueOnce()
	}
	ctx.Data["Title"] = ctx.Tr("admin.monitor.queues")
	ctx.Data["PageIsAdminMonitorQueue"] = true
	ctx.Data["Queues"] = queue.GetManager().ManagedQueues()
	ctx.HTML(http.StatusOK, tplQueue)
}

// QueueManage shows details for a specific queue
func QueueManage(ctx *context.Context) {
	qid := ctx.PathParamInt64("qid")
	mq := queue.GetManager().GetManagedQueue(qid)
	if mq == nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.Data["Title"] = ctx.Tr("admin.monitor.queue", mq.GetName())
	ctx.Data["PageIsAdminMonitor"] = true
	ctx.Data["Queue"] = mq
	ctx.HTML(http.StatusOK, tplQueueManage)
}

// QueueSet sets the maximum number of workers and other settings for this queue
func QueueSet(ctx *context.Context) {
	qid := ctx.PathParamInt64("qid")
	mq := queue.GetManager().GetManagedQueue(qid)
	if mq == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	maxNumberStr := ctx.FormString("max-number")

	var err error
	var maxNumber int
	if len(maxNumberStr) > 0 {
		maxNumber, err = strconv.Atoi(maxNumberStr)
		if err != nil {
			ctx.Flash.Error(ctx.Tr("admin.monitor.queue.settings.maxnumberworkers.error"))
			ctx.Redirect(setting.AppSubURL + "/-/admin/monitor/queue/" + strconv.FormatInt(qid, 10))
			return
		}
		if maxNumber < -1 {
			maxNumber = -1
		}
	} else {
		maxNumber = mq.GetWorkerMaxNumber()
	}

	mq.SetWorkerMaxNumber(maxNumber)
	ctx.Flash.Success(ctx.Tr("admin.monitor.queue.settings.changed"))
	ctx.Redirect(setting.AppSubURL + "/-/admin/monitor/queue/" + strconv.FormatInt(qid, 10))
}

func QueueRemoveAllItems(ctx *context.Context) {
	// Kmup's queue doesn't have transaction support
	// So in rare cases, the queue could be corrupted/out-of-sync
	// Site admin could remove all items from the queue to make it work again
	qid := ctx.PathParamInt64("qid")
	mq := queue.GetManager().GetManagedQueue(qid)
	if mq == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	if err := mq.RemoveAllItems(ctx); err != nil {
		ctx.ServerError("RemoveAllItems", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("admin.monitor.queue.settings.remove_all_items_done"))
	ctx.Redirect(setting.AppSubURL + "/-/admin/monitor/queue/" + strconv.FormatInt(qid, 10))
}
