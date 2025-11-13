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

package v1_20

import (
	"context"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/setting"

	"xorm.io/xorm"
)

func RenameWebhookOrgToOwner(x *xorm.Engine) error {
	type Webhook struct {
		OrgID int64 `xorm:"INDEX"`
	}

	// This migration maybe rerun so that we should check if it has been run
	ownerExist, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "webhook", "owner_id")
	if err != nil {
		return err
	}

	if ownerExist {
		orgExist, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "webhook", "org_id")
		if err != nil {
			return err
		}
		if !orgExist {
			return nil
		}
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := sess.Sync(new(Webhook)); err != nil {
		return err
	}

	if ownerExist {
		if err := base.DropTableColumns(sess, "webhook", "owner_id"); err != nil {
			return err
		}
	}

	switch {
	case setting.Database.Type.IsMySQL():
		inferredTable, err := x.TableInfo(new(Webhook))
		if err != nil {
			return err
		}
		sqlType := x.Dialect().SQLType(inferredTable.GetColumn("org_id"))
		if _, err := sess.Exec("ALTER TABLE `webhook` CHANGE org_id owner_id " + sqlType); err != nil {
			return err
		}
	case setting.Database.Type.IsMSSQL():
		if _, err := sess.Exec("sp_rename 'webhook.org_id', 'owner_id', 'COLUMN'"); err != nil {
			return err
		}
	default:
		if _, err := sess.Exec("ALTER TABLE `webhook` RENAME COLUMN org_id TO owner_id"); err != nil {
			return err
		}
	}

	return sess.Commit()
}
