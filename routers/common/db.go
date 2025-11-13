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

package common

import (
	"context"
	"errors"
	"time"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/migrations"
	system_model "github.com/kumose/kmup/models/system"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/setting/config"
	"github.com/kumose/kmup/services/versioned_migration"

	"xorm.io/xorm"
)

// InitDBEngine In case of problems connecting to DB, retry connection. Eg, PGSQL in Docker Container on Synology
func InitDBEngine(ctx context.Context) (err error) {
	log.Info("Beginning ORM engine initialization.")
	for i := 0; i < setting.Database.DBConnectRetries; i++ {
		select {
		case <-ctx.Done():
			return errors.New("Aborted due to shutdown:\nin retry ORM engine initialization")
		default:
		}
		log.Info("ORM engine initialization attempt #%d/%d...", i+1, setting.Database.DBConnectRetries)
		if err = db.InitEngineWithMigration(ctx, migrateWithSetting); err == nil {
			break
		} else if i == setting.Database.DBConnectRetries-1 {
			return err
		}
		log.Error("ORM engine initialization attempt #%d/%d failed. Error: %v", i+1, setting.Database.DBConnectRetries, err)
		log.Info("Backing off for %d seconds", int64(setting.Database.DBConnectBackoff/time.Second))
		time.Sleep(setting.Database.DBConnectBackoff)
	}
	config.SetDynGetter(system_model.NewDatabaseDynKeyGetter())
	return nil
}

func migrateWithSetting(ctx context.Context, x *xorm.Engine) error {
	if setting.Database.AutoMigration {
		return versioned_migration.Migrate(ctx, x)
	}

	if current, err := migrations.GetCurrentDBVersion(x); err != nil {
		return err
	} else if current < 0 {
		// execute migrations when the database isn't initialized even if AutoMigration is false
		return versioned_migration.Migrate(ctx, x)
	} else if expected := migrations.ExpectedDBVersion(); current != expected {
		log.Fatal(`"database.AUTO_MIGRATION" is disabled, but current database version %d is not equal to the expected version %d.`+
			`You can set "database.AUTO_MIGRATION" to true or migrate manually by running "kmup [--config /path/to/app.ini] migrate"`, current, expected)
	}
	return nil
}
