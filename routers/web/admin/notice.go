// Copyright 2014 The Gogs Authors. All rights reserved.
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

	"github.com/kumose/kmup/models/db"
	system_model "github.com/kumose/kmup/models/system"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/services/context"
)

const (
	tplNotices templates.TplName = "admin/notice"
)

// Notices show notices for admin
func Notices(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.notices")
	ctx.Data["PageIsAdminNotices"] = true

	total := system_model.CountNotices(ctx)
	page := max(ctx.FormInt("page"), 1)

	notices, err := system_model.Notices(ctx, page, setting.UI.Admin.NoticePagingNum)
	if err != nil {
		ctx.ServerError("Notices", err)
		return
	}
	ctx.Data["Notices"] = notices

	ctx.Data["Total"] = total

	ctx.Data["Page"] = context.NewPagination(int(total), setting.UI.Admin.NoticePagingNum, page, 5)

	ctx.HTML(http.StatusOK, tplNotices)
}

// DeleteNotices delete the specific notices
func DeleteNotices(ctx *context.Context) {
	strs := ctx.FormStrings("ids[]")
	ids := make([]int64, 0, len(strs))
	for i := range strs {
		id, _ := strconv.ParseInt(strs[i], 10, 64)
		if id > 0 {
			ids = append(ids, id)
		}
	}

	if err := db.DeleteByIDs[system_model.Notice](ctx, ids...); err != nil {
		ctx.Flash.Error("DeleteNoticesByIDs: " + err.Error())
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.Flash.Success(ctx.Tr("admin.notices.delete_success"))
		ctx.Status(http.StatusOK)
	}
}

// EmptyNotices delete all the notices
func EmptyNotices(ctx *context.Context) {
	if err := system_model.DeleteNotices(ctx, 0, 0); err != nil {
		ctx.ServerError("DeleteNotices", err)
		return
	}

	log.Trace("System notices deleted by admin (%s): [start: %d]", ctx.Doer.Name, 0)
	ctx.Flash.Success(ctx.Tr("admin.notices.delete_success"))
	ctx.Redirect(setting.AppSubURL + "/-/admin/notices")
}
