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
	"context"

	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/modules/util"
	secret_service "github.com/kumose/kmup/services/secrets"
)

func CreateVariable(ctx context.Context, ownerID, repoID int64, name, data, description string) (*actions_model.ActionVariable, error) {
	if err := secret_service.ValidateName(name); err != nil {
		return nil, err
	}

	v, err := actions_model.InsertVariable(ctx, ownerID, repoID, name, util.ReserveLineBreakForTextarea(data), description)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func UpdateVariableNameData(ctx context.Context, variable *actions_model.ActionVariable) (bool, error) {
	if err := secret_service.ValidateName(variable.Name); err != nil {
		return false, err
	}

	variable.Data = util.ReserveLineBreakForTextarea(variable.Data)

	return actions_model.UpdateVariableCols(ctx, variable, "name", "data", "description")
}

func DeleteVariableByID(ctx context.Context, variableID int64) error {
	return actions_model.DeleteVariable(ctx, variableID)
}

func DeleteVariableByName(ctx context.Context, ownerID, repoID int64, name string) error {
	v, err := GetVariable(ctx, actions_model.FindVariablesOpts{
		OwnerID: ownerID,
		RepoID:  repoID,
		Name:    name,
	})
	if err != nil {
		return err
	}

	return actions_model.DeleteVariable(ctx, v.ID)
}

func GetVariable(ctx context.Context, opts actions_model.FindVariablesOpts) (*actions_model.ActionVariable, error) {
	vars, err := actions_model.FindVariables(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(vars) != 1 {
		return nil, util.NewNotExistErrorf("variable not found")
	}
	return vars[0], nil
}
