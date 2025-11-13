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

package notify

import (
	"fmt"
	"net/http"

	activities_model "github.com/kumose/kmup/models/activities"
	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// GetThread get notification by ID
func GetThread(ctx *context.APIContext) {
	// swagger:operation GET /notifications/threads/{id} notification notifyGetThread
	// ---
	// summary: Get notification thread by ID
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of notification thread
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/NotificationThread"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	//   "404":
	//     "$ref": "#/responses/notFound"

	n := getThread(ctx)
	if n == nil {
		return
	}
	if err := n.LoadAttributes(ctx); err != nil && !issues_model.IsErrCommentNotExist(err) {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.JSON(http.StatusOK, convert.ToNotificationThread(ctx, n))
}

// ReadThread mark notification as read by ID
func ReadThread(ctx *context.APIContext) {
	// swagger:operation PATCH /notifications/threads/{id} notification notifyReadThread
	// ---
	// summary: Mark notification thread as read by ID
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of notification thread
	//   type: string
	//   required: true
	// - name: to-status
	//   in: query
	//   description: Status to mark notifications as
	//   type: string
	//   default: read
	//   required: false
	// responses:
	//   "205":
	//     "$ref": "#/responses/NotificationThread"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	//   "404":
	//     "$ref": "#/responses/notFound"

	n := getThread(ctx)
	if n == nil {
		return
	}

	targetStatus := statusStringToNotificationStatus(ctx.FormString("to-status"))
	if targetStatus == 0 {
		targetStatus = activities_model.NotificationStatusRead
	}

	notif, err := activities_model.SetNotificationStatus(ctx, n.ID, ctx.Doer, targetStatus)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	if err = notif.LoadAttributes(ctx); err != nil && !issues_model.IsErrCommentNotExist(err) {
		ctx.APIErrorInternal(err)
		return
	}
	ctx.JSON(http.StatusResetContent, convert.ToNotificationThread(ctx, notif))
}

func getThread(ctx *context.APIContext) *activities_model.Notification {
	n, err := activities_model.GetNotificationByID(ctx, ctx.PathParamInt64("id"))
	if err != nil {
		if db.IsErrNotExist(err) {
			ctx.APIError(http.StatusNotFound, err)
		} else {
			ctx.APIErrorInternal(err)
		}
		return nil
	}
	if n.UserID != ctx.Doer.ID && !ctx.Doer.IsAdmin {
		ctx.APIError(http.StatusForbidden, fmt.Errorf("only user itself and admin are allowed to read/change this thread %d", n.ID))
		return nil
	}
	return n
}
