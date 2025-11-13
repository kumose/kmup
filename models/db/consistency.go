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

package db

import (
	"context"

	"xorm.io/builder"
)

// CountOrphanedObjects count subjects with have no existing refobject anymore
func CountOrphanedObjects(ctx context.Context, subject, refObject, joinCond string) (int64, error) {
	return GetEngine(ctx).
		Table("`"+subject+"`").
		Join("LEFT", "`"+refObject+"`", joinCond).
		Where(builder.IsNull{"`" + refObject + "`.id"}).
		Select("COUNT(`" + subject + "`.`id`)").
		Count()
}

// DeleteOrphanedObjects delete subjects with have no existing refobject anymore
func DeleteOrphanedObjects(ctx context.Context, subject, refObject, joinCond string) error {
	subQuery := builder.Select("`"+subject+"`.id").
		From("`"+subject+"`").
		Join("LEFT", "`"+refObject+"`", joinCond).
		Where(builder.IsNull{"`" + refObject + "`.id"})
	b := builder.Delete(builder.In("id", subQuery)).From("`" + subject + "`")
	_, err := GetEngine(ctx).Exec(b)
	return err
}
