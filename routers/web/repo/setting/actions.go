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

package setting

import (
	"errors"
	"net/http"
	"strings"

	repo_model "github.com/kumose/kmup/models/repo"
	unit_model "github.com/kumose/kmup/models/unit"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
	repo_service "github.com/kumose/kmup/services/repository"
)

const tplRepoActionsGeneralSettings templates.TplName = "repo/settings/actions"

func ActionsGeneralSettings(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("actions.general")
	ctx.Data["PageType"] = "general"
	ctx.Data["PageIsActionsSettingsGeneral"] = true

	actionsUnit, err := ctx.Repo.Repository.GetUnit(ctx, unit_model.TypeActions)
	if err != nil && !repo_model.IsErrUnitTypeNotExist(err) {
		ctx.ServerError("GetUnit", err)
		return
	}
	if actionsUnit == nil { // no actions unit
		ctx.HTML(http.StatusOK, tplRepoActionsGeneralSettings)
		return
	}

	if ctx.Repo.Repository.IsPrivate {
		collaborativeOwnerIDs := actionsUnit.ActionsConfig().CollaborativeOwnerIDs
		collaborativeOwners, err := user_model.GetUsersByIDs(ctx, collaborativeOwnerIDs)
		if err != nil {
			ctx.ServerError("GetUsersByIDs", err)
			return
		}
		ctx.Data["CollaborativeOwners"] = collaborativeOwners
	}

	ctx.HTML(http.StatusOK, tplRepoActionsGeneralSettings)
}

func ActionsUnitPost(ctx *context.Context) {
	redirectURL := ctx.Repo.RepoLink + "/settings/actions/general"
	enableActionsUnit := ctx.FormBool("enable_actions")
	repo := ctx.Repo.Repository

	var err error
	if enableActionsUnit && !unit_model.TypeActions.UnitGlobalDisabled() {
		err = repo_service.UpdateRepositoryUnits(ctx, repo, []repo_model.RepoUnit{newRepoUnit(repo, unit_model.TypeActions, nil)}, nil)
	} else if !unit_model.TypeActions.UnitGlobalDisabled() {
		err = repo_service.UpdateRepositoryUnits(ctx, repo, nil, []unit_model.Type{unit_model.TypeActions})
	}

	if err != nil {
		ctx.ServerError("UpdateRepositoryUnits", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("repo.settings.update_settings_success"))
	ctx.Redirect(redirectURL)
}

func AddCollaborativeOwner(ctx *context.Context) {
	name := strings.ToLower(ctx.FormString("collaborative_owner"))

	ownerID, err := user_model.GetUserOrOrgIDByName(ctx, name)
	if err != nil {
		if errors.Is(err, util.ErrNotExist) {
			ctx.Flash.Error(ctx.Tr("form.user_not_exist"))
			ctx.JSONErrorNotFound()
		} else {
			ctx.ServerError("GetUserOrOrgIDByName", err)
		}
		return
	}

	actionsUnit, err := ctx.Repo.Repository.GetUnit(ctx, unit_model.TypeActions)
	if err != nil {
		ctx.ServerError("GetUnit", err)
		return
	}
	actionsCfg := actionsUnit.ActionsConfig()
	actionsCfg.AddCollaborativeOwner(ownerID)
	if err := repo_model.UpdateRepoUnit(ctx, actionsUnit); err != nil {
		ctx.ServerError("UpdateRepoUnit", err)
		return
	}

	ctx.JSONOK()
}

func DeleteCollaborativeOwner(ctx *context.Context) {
	ownerID := ctx.FormInt64("id")

	actionsUnit, err := ctx.Repo.Repository.GetUnit(ctx, unit_model.TypeActions)
	if err != nil {
		ctx.ServerError("GetUnit", err)
		return
	}
	actionsCfg := actionsUnit.ActionsConfig()
	if !actionsCfg.IsCollaborativeOwner(ownerID) {
		ctx.Flash.Error(ctx.Tr("actions.general.collaborative_owner_not_exist"))
		ctx.JSONErrorNotFound()
		return
	}
	actionsCfg.RemoveCollaborativeOwner(ownerID)
	if err := repo_model.UpdateRepoUnit(ctx, actionsUnit); err != nil {
		ctx.ServerError("UpdateRepoUnit", err)
		return
	}

	ctx.JSONOK()
}
