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

package actions

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/modules/badge"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
)

func GetWorkflowBadge(ctx *context.Context) {
	workflowFile := ctx.PathParam("workflow_name")
	branch := ctx.FormString("branch", ctx.Repo.Repository.DefaultBranch)
	event := ctx.FormString("event")
	style := ctx.FormString("style")

	branchRef := git.RefNameFromBranch(branch)
	b, err := getWorkflowBadge(ctx, workflowFile, branchRef.String(), event)
	if err != nil {
		ctx.ServerError("GetWorkflowBadge", err)
		return
	}

	ctx.Data["Badge"] = b
	ctx.RespHeader().Set("Content-Type", "image/svg+xml")
	switch style {
	case badge.StyleFlatSquare:
		ctx.HTML(http.StatusOK, "shared/actions/runner_badge_flat-square")
	default: // defaults to badge.StyleFlat
		ctx.HTML(http.StatusOK, "shared/actions/runner_badge_flat")
	}
}

func getWorkflowBadge(ctx *context.Context, workflowFile, branchName, event string) (badge.Badge, error) {
	extension := filepath.Ext(workflowFile)
	workflowName := strings.TrimSuffix(workflowFile, extension)

	run, err := actions_model.GetWorkflowLatestRun(ctx, ctx.Repo.Repository.ID, workflowFile, branchName, event)
	if err != nil {
		if errors.Is(err, util.ErrNotExist) {
			return badge.GenerateBadge(workflowName, "no status", badge.DefaultColor), nil
		}
		return badge.Badge{}, err
	}

	color, ok := badge.GlobalVars().StatusColorMap[run.Status]
	if !ok {
		return badge.GenerateBadge(workflowName, "unknown status", badge.DefaultColor), nil
	}
	return badge.GenerateBadge(workflowName, run.Status.String(), color), nil
}
