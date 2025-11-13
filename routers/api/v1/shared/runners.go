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

package shared

import (
	"errors"
	"net/http"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/routers/api/v1/utils"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/convert"
)

// RegistrationToken is response related to registration token
// swagger:response RegistrationToken
type RegistrationToken struct {
	Token string `json:"token"`
}

func GetRegistrationToken(ctx *context.APIContext, ownerID, repoID int64) {
	token, err := actions_model.GetLatestRunnerToken(ctx, ownerID, repoID)
	if errors.Is(err, util.ErrNotExist) || (token != nil && !token.IsActive) {
		token, err = actions_model.NewRunnerToken(ctx, ownerID, repoID)
	}
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	ctx.JSON(http.StatusOK, RegistrationToken{Token: token.Token})
}

// ListRunners lists runners for api route validated ownerID and repoID
// ownerID == 0 and repoID == 0 means all runners including global runners, does not appear in sql where clause
// ownerID == 0 and repoID != 0 means all runners for the given repo
// ownerID != 0 and repoID == 0 means all runners for the given user/org
// ownerID != 0 and repoID != 0 undefined behavior
// Access rights are checked at the API route level
func ListRunners(ctx *context.APIContext, ownerID, repoID int64) {
	if ownerID != 0 && repoID != 0 {
		setting.PanicInDevOrTesting("ownerID and repoID should not be both set")
	}
	runners, total, err := db.FindAndCount[actions_model.ActionRunner](ctx, &actions_model.FindRunnerOptions{
		OwnerID:     ownerID,
		RepoID:      repoID,
		ListOptions: utils.GetListOptions(ctx),
	})
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	res := new(api.ActionRunnersResponse)
	res.TotalCount = total

	res.Entries = make([]*api.ActionRunner, len(runners))
	for i, runner := range runners {
		res.Entries[i] = convert.ToActionRunner(ctx, runner)
	}

	ctx.JSON(http.StatusOK, &res)
}

func getRunnerByID(ctx *context.APIContext, ownerID, repoID, runnerID int64) (*actions_model.ActionRunner, bool) {
	if ownerID != 0 && repoID != 0 {
		setting.PanicInDevOrTesting("ownerID and repoID should not be both set")
	}

	runner, err := actions_model.GetRunnerByID(ctx, runnerID)
	if err != nil {
		if errors.Is(err, util.ErrNotExist) {
			ctx.APIErrorNotFound("Runner not found")
		} else {
			ctx.APIErrorInternal(err)
		}
		return nil, false
	}

	if !runner.EditableInContext(ownerID, repoID) {
		ctx.APIErrorNotFound("No permission to access this runner")
		return nil, false
	}
	return runner, true
}

// GetRunner get the runner for api route validated ownerID and repoID
// ownerID == 0 and repoID == 0 means any runner including global runners
// ownerID == 0 and repoID != 0 means any runner for the given repo
// ownerID != 0 and repoID == 0 means any runner for the given user/org
// ownerID != 0 and repoID != 0 undefined behavior
// Access rights are checked at the API route level
func GetRunner(ctx *context.APIContext, ownerID, repoID, runnerID int64) {
	if ownerID != 0 && repoID != 0 {
		setting.PanicInDevOrTesting("ownerID and repoID should not be both set")
	}
	runner, ok := getRunnerByID(ctx, ownerID, repoID, runnerID)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, convert.ToActionRunner(ctx, runner))
}

// DeleteRunner deletes the runner for api route validated ownerID and repoID
// ownerID == 0 and repoID == 0 means any runner including global runners
// ownerID == 0 and repoID != 0 means any runner for the given repo
// ownerID != 0 and repoID == 0 means any runner for the given user/org
// ownerID != 0 and repoID != 0 undefined behavior
// Access rights are checked at the API route level
func DeleteRunner(ctx *context.APIContext, ownerID, repoID, runnerID int64) {
	runner, ok := getRunnerByID(ctx, ownerID, repoID, runnerID)
	if !ok {
		return
	}

	err := actions_model.DeleteRunner(ctx, runner.ID)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
