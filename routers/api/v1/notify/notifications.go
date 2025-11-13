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
	"net/http"
	"strings"

	activities_model "github.com/kumose/kmup/models/activities"
	"github.com/kumose/kmup/models/db"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
)

// NewAvailable check if unread notifications exist
func NewAvailable(ctx *context.APIContext) {
	// swagger:operation GET /notifications/new notification notifyNewAvailable
	// ---
	// summary: Check if unread notifications exist
	// responses:
	//   "200":
	//     "$ref": "#/responses/NotificationCount"

	total, err := db.Count[activities_model.Notification](ctx, activities_model.FindNotificationOptions{
		UserID: ctx.Doer.ID,
		Status: []activities_model.NotificationStatus{activities_model.NotificationStatusUnread},
	})
	if err != nil {
		ctx.APIError(http.StatusUnprocessableEntity, err)
		return
	}

	ctx.JSON(http.StatusOK, api.NotificationCount{New: total})
}

func getFindNotificationOptions(ctx *context.APIContext) *activities_model.FindNotificationOptions {
	before, since, err := context.GetQueryBeforeSince(ctx.Base)
	if err != nil {
		ctx.APIError(http.StatusUnprocessableEntity, err)
		return nil
	}
	opts := &activities_model.FindNotificationOptions{
		ListOptions:       utils.GetListOptions(ctx),
		UserID:            ctx.Doer.ID,
		UpdatedBeforeUnix: before,
		UpdatedAfterUnix:  since,
	}
	if !ctx.FormBool("all") {
		statuses := ctx.FormStrings("status-types")
		opts.Status = statusStringsToNotificationStatuses(statuses, []string{"unread", "pinned"})
	}

	subjectTypes := ctx.FormStrings("subject-type")
	if len(subjectTypes) != 0 {
		opts.Source = subjectToSource(subjectTypes)
	}

	return opts
}

func subjectToSource(value []string) (result []activities_model.NotificationSource) {
	for _, v := range value {
		switch strings.ToLower(v) {
		case "issue":
			result = append(result, activities_model.NotificationSourceIssue)
		case "pull":
			result = append(result, activities_model.NotificationSourcePullRequest)
		case "commit":
			result = append(result, activities_model.NotificationSourceCommit)
		case "repository":
			result = append(result, activities_model.NotificationSourceRepository)
		}
	}
	return result
}
