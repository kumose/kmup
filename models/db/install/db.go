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

package install

import (
	"context"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/setting"
)

// CheckDatabaseConnection checks the database connection
func CheckDatabaseConnection(ctx context.Context) error {
	_, err := db.GetEngine(ctx).Exec("SELECT 1")
	return err
}

// GetMigrationVersion gets the database migration version
func GetMigrationVersion(ctx context.Context) (int64, error) {
	var installedDbVersion int64
	x := db.GetEngine(ctx)
	exist, err := x.IsTableExist("version")
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, nil
	}
	_, err = x.Table("version").Cols("version").Get(&installedDbVersion)
	if err != nil {
		return 0, err
	}
	return installedDbVersion, nil
}

// HasPostInstallationUsers checks whether there are users after installation
func HasPostInstallationUsers(ctx context.Context) (bool, error) {
	x := db.GetEngine(ctx)
	exist, err := x.IsTableExist("user")
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}

	// if there are 2 or more users in database, we consider there are users created after installation
	threshold := 2
	if !setting.IsProd {
		// to debug easily, with non-prod RUN_MODE, we only check the count to 1
		threshold = 1
	}
	res, err := x.Table("user").Cols("id").Limit(threshold).Query()
	if err != nil {
		return false, err
	}
	return len(res) >= threshold, nil
}
