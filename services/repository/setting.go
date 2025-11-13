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

package repository

import (
	"context"
	"slices"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/modules/log"
	actions_service "github.com/kumose/kmup/services/actions"
)

// UpdateRepositoryUnits updates a repository's units
func UpdateRepositoryUnits(ctx context.Context, repo *repo_model.Repository, units []repo_model.RepoUnit, deleteUnitTypes []unit.Type) (err error) {
	return db.WithTx(ctx, func(ctx context.Context) error {
		// Delete existing settings of units before adding again
		for _, u := range units {
			deleteUnitTypes = append(deleteUnitTypes, u.Type)
		}

		if slices.Contains(deleteUnitTypes, unit.TypeActions) {
			if err := actions_service.CleanRepoScheduleTasks(ctx, repo); err != nil {
				log.Error("CleanRepoScheduleTasks: %v", err)
			}
		}

		for _, u := range units {
			if u.Type == unit.TypeActions {
				if err := actions_service.DetectAndHandleSchedules(ctx, repo); err != nil {
					log.Error("DetectAndHandleSchedules: %v", err)
				}
				break
			}
		}

		if _, err = db.GetEngine(ctx).Where("repo_id = ?", repo.ID).In("type", deleteUnitTypes).Delete(new(repo_model.RepoUnit)); err != nil {
			return err
		}

		if len(units) > 0 {
			if err = db.Insert(ctx, units); err != nil {
				return err
			}
		}

		return nil
	})
}
