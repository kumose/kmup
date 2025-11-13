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

	"github.com/kumose/kmup/models/organization"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
)

// CanUserForkBetweenOwners returns true if user can fork between owners.
// By default, a user can fork a repository from another owner, but not from themselves.
// Many users really like to fork their own repositories, so add an experimental setting to allow this.
func CanUserForkBetweenOwners(id1, id2 int64) bool {
	if id1 != id2 {
		return true
	}
	return setting.Repository.AllowForkIntoSameOwner
}

// CanUserForkRepo returns true if specified user can fork repository.
func CanUserForkRepo(ctx context.Context, user *user_model.User, repo *repo_model.Repository) (bool, error) {
	if user == nil {
		return false, nil
	}
	if CanUserForkBetweenOwners(repo.OwnerID, user.ID) && !repo_model.HasForkedRepo(ctx, user.ID, repo.ID) {
		return true, nil
	}
	ownedOrgs, err := organization.GetOrgsCanCreateRepoByUserID(ctx, user.ID)
	if err != nil {
		return false, err
	}
	for _, org := range ownedOrgs {
		if repo.OwnerID != org.ID && !repo_model.HasForkedRepo(ctx, org.ID, repo.ID) {
			return true, nil
		}
	}
	return false, nil
}
