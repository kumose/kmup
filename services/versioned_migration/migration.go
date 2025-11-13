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

package versioned_migration

import (
	"context"

	"github.com/kumose/kmup/models/migrations"
	"github.com/kumose/kmup/modules/globallock"

	"xorm.io/xorm"
)

func Migrate(ctx context.Context, x *xorm.Engine) error {
	// only one instance can do the migration at the same time if there are multiple instances
	release, err := globallock.Lock(ctx, "kmup_versioned_migration")
	if err != nil {
		return err
	}
	defer release()

	return migrations.Migrate(ctx, x)
}
