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

package migrations

import (
	"context"

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/services/externalaccount"
)

// UpdateMigrationPosterID updates all migrated repositories' issues and comments posterID
func UpdateMigrationPosterID(ctx context.Context) error {
	for _, gitService := range structs.SupportedFullGitService {
		select {
		case <-ctx.Done():
			log.Warn("UpdateMigrationPosterID aborted before %s", gitService.Name())
			return db.ErrCancelledf("during UpdateMigrationPosterID before %s", gitService.Name())
		default:
		}
		if err := updateMigrationPosterIDByGitService(ctx, gitService); err != nil {
			log.Error("updateMigrationPosterIDByGitService failed: %v", err)
		}
	}
	return nil
}

func updateMigrationPosterIDByGitService(ctx context.Context, tp structs.GitServiceType) error {
	provider := tp.Name()
	if len(provider) == 0 {
		return nil
	}

	const batchSize = 100
	for page := 0; ; page++ {
		select {
		case <-ctx.Done():
			log.Warn("UpdateMigrationPosterIDByGitService(%s) cancelled", tp.Name())
			return nil
		default:
		}

		users, err := db.Find[user_model.ExternalLoginUser](ctx, user_model.FindExternalUserOptions{
			ListOptions: db.ListOptions{
				PageSize: batchSize,
				Page:     page,
			},
			Provider: provider,
			OrderBy:  "login_source_id ASC, external_id ASC",
		})
		if err != nil {
			return err
		}

		for _, user := range users {
			select {
			case <-ctx.Done():
				log.Warn("UpdateMigrationPosterIDByGitService(%s) cancelled", tp.Name())
				return nil
			default:
			}
			externalUserID := user.ExternalID
			if err := externalaccount.UpdateMigrationsByType(ctx, tp, externalUserID, user.UserID); err != nil {
				log.Error("UpdateMigrationsByType type %s external user id %v to local user id %v failed: %v", tp.Name(), user.ExternalID, user.UserID, err)
			}
		}

		if len(users) < batchSize {
			break
		}
	}
	return nil
}
