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

package user

import (
	"net/http"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// GetStopwatches get all stopwatches
func GetStopwatches(ctx *context.Context) {
	sws, err := issues_model.GetUserStopwatches(ctx, ctx.Doer.ID, db.ListOptions{
		Page:     ctx.FormInt("page"),
		PageSize: convert.ToCorrectPageSize(ctx.FormInt("limit")),
	})
	if err != nil {
		ctx.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}

	count, err := issues_model.CountUserStopwatches(ctx, ctx.Doer.ID)
	if err != nil {
		ctx.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}

	apiSWs, err := convert.ToStopWatches(ctx, sws)
	if err != nil {
		ctx.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetTotalCountHeader(count)
	ctx.JSON(http.StatusOK, apiSWs)
}
