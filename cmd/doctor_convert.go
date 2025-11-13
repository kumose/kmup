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
	"fmt"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"github.com/urfave/cli/v3"
)

// cmdDoctorConvert represents the available convert sub-command.
var cmdDoctorConvert = &cli.Command{
	Name:        "convert",
	Usage:       "Convert the database",
	Description: "A command to convert an existing MySQL database from utf8 to utf8mb4 or MSSQL database from varchar to nvarchar",
	Action:      runDoctorConvert,
}

func runDoctorConvert(ctx context.Context, cmd *cli.Command) error {
	if err := initDB(ctx); err != nil {
		return err
	}

	log.Info("AppPath: %s", setting.AppPath)
	log.Info("AppWorkPath: %s", setting.AppWorkPath)
	log.Info("Custom path: %s", setting.CustomPath)
	log.Info("Log path: %s", setting.Log.RootPath)
	log.Info("Configuration file: %s", setting.CustomConf)

	switch {
	case setting.Database.Type.IsMySQL():
		if err := db.ConvertDatabaseTable(); err != nil {
			log.Fatal("Failed to convert database & table: %v", err)
			return err
		}
		fmt.Println("Converted successfully, please confirm your database's character set is now utf8mb4")
	case setting.Database.Type.IsMSSQL():
		if err := db.ConvertVarcharToNVarchar(); err != nil {
			log.Fatal("Failed to convert database from varchar to nvarchar: %v", err)
			return err
		}
		fmt.Println("Converted successfully, please confirm your database's all columns character is NVARCHAR now")
	default:
		fmt.Println("This command can only be used with a MySQL or MSSQL database")
	}

	return nil
}
