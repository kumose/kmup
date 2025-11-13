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

package cmd

import (
	"context"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/versioned_migration"

	"github.com/urfave/cli/v3"
)

// CmdMigrate represents the available migrate sub-command.
var CmdMigrate = &cli.Command{
	Name:        "migrate",
	Usage:       "Migrate the database",
	Description: `This is a command for migrating the database, so that you can run "kmup admin create user" before starting the server.`,
	Action:      runMigrate,
}

func runMigrate(ctx context.Context, c *cli.Command) error {
	if err := initDB(ctx); err != nil {
		return err
	}

	log.Info("AppPath: %s", setting.AppPath)
	log.Info("AppWorkPath: %s", setting.AppWorkPath)
	log.Info("Custom path: %s", setting.CustomPath)
	log.Info("Log path: %s", setting.Log.RootPath)
	log.Info("Configuration file: %s", setting.CustomConf)

	if err := db.InitEngineWithMigration(context.Background(), versioned_migration.Migrate); err != nil {
		log.Fatal("Failed to initialize ORM engine: %v", err)
		return err
	}

	return nil
}
