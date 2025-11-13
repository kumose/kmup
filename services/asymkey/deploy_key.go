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

package asymkey

import (
	"context"
	"fmt"

	asymkey_model "github.com/kumose/kmup/models/asymkey"
	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
)

// DeleteRepoDeployKeys deletes all deploy keys of a repository. permissions check should be done outside
func DeleteRepoDeployKeys(ctx context.Context, repoID int64) (int, error) {
	deployKeys, err := db.Find[asymkey_model.DeployKey](ctx, asymkey_model.ListDeployKeysOptions{RepoID: repoID})
	if err != nil {
		return 0, fmt.Errorf("listDeployKeys: %w", err)
	}

	for _, dKey := range deployKeys {
		if err := deleteDeployKeyFromDB(ctx, dKey); err != nil {
			return 0, fmt.Errorf("deleteDeployKeys: %w", err)
		}
	}
	return len(deployKeys), nil
}

// deleteDeployKeyFromDB delete deploy keys from database
func deleteDeployKeyFromDB(ctx context.Context, key *asymkey_model.DeployKey) error {
	if _, err := db.DeleteByID[asymkey_model.DeployKey](ctx, key.ID); err != nil {
		return fmt.Errorf("delete deploy key [%d]: %w", key.ID, err)
	}

	// Check if this is the last reference to same key content.
	has, err := asymkey_model.IsDeployKeyExistByKeyID(ctx, key.KeyID)
	if err != nil {
		return err
	} else if !has {
		if _, err = db.DeleteByID[asymkey_model.PublicKey](ctx, key.KeyID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteDeployKey deletes deploy key from its repository authorized_keys file if needed.
// Permissions check should be done outside.
func DeleteDeployKey(ctx context.Context, repo *repo_model.Repository, id int64) error {
	if err := db.WithTx(ctx, func(ctx context.Context) error {
		key, err := asymkey_model.GetDeployKeyByID(ctx, id)
		if err != nil {
			if asymkey_model.IsErrDeployKeyNotExist(err) {
				return nil
			}
			return fmt.Errorf("GetDeployKeyByID: %w", err)
		}

		if key.RepoID != repo.ID {
			return fmt.Errorf("deploy key %d does not belong to repository %d", id, repo.ID)
		}

		return deleteDeployKeyFromDB(ctx, key)
	}); err != nil {
		return err
	}

	return RewriteAllPublicKeys(ctx)
}
