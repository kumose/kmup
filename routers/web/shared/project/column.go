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

package project

import (
	project_model "github.com/kumose/kmup/models/project"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/services/context"
)

// MoveColumns moves or keeps columns in a project and sorts them inside that project
func MoveColumns(ctx *context.Context) {
	project, err := project_model.GetProjectByID(ctx, ctx.PathParamInt64("id"))
	if err != nil {
		ctx.NotFoundOrServerError("GetProjectByID", project_model.IsErrProjectNotExist, err)
		return
	}
	if !project.CanBeAccessedByOwnerRepo(ctx.ContextUser.ID, ctx.Repo.Repository) {
		ctx.NotFound(nil)
		return
	}

	type movedColumnsForm struct {
		Columns []struct {
			ColumnID int64 `json:"columnID"`
			Sorting  int64 `json:"sorting"`
		} `json:"columns"`
	}

	form := &movedColumnsForm{}
	if err = json.NewDecoder(ctx.Req.Body).Decode(&form); err != nil {
		ctx.ServerError("DecodeMovedColumnsForm", err)
		return
	}

	sortedColumnIDs := make(map[int64]int64)
	for _, column := range form.Columns {
		sortedColumnIDs[column.Sorting] = column.ColumnID
	}

	if err = project_model.MoveColumnsOnProject(ctx, project, sortedColumnIDs); err != nil {
		ctx.ServerError("MoveColumnsOnProject", err)
		return
	}

	ctx.JSONOK()
}
