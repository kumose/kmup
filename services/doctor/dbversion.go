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

package doctor

import (
	"context"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/migrations"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/versioned_migration"
)

func checkDBVersion(ctx context.Context, logger log.Logger, autofix bool) error {
	logger.Info("Expected database version: %d", migrations.ExpectedDBVersion())
	if err := db.InitEngineWithMigration(ctx, migrations.EnsureUpToDate); err != nil {
		if !autofix {
			logger.Critical("Error: %v during ensure up to date", err)
			return err
		}
		logger.Warn("Got Error: %v during ensure up to date", err)
		logger.Warn("Attempting to migrate to the latest DB version to fix this.")

		err = db.InitEngineWithMigration(ctx, versioned_migration.Migrate)
		if err != nil {
			logger.Critical("Error: %v during migration", err)
		}
		return err
	}
	return nil
}

func init() {
	Register(&Check{
		Title:         "Check Database Version",
		Name:          "check-db-version",
		IsDefault:     true,
		Run:           checkDBVersion,
		AbortIfFailed: false,
		Priority:      2,
	})
}
