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

package v1_12

import (
	"xorm.io/builder"
	"xorm.io/xorm"
)

func FixTopicRepositoryCount(x *xorm.Engine) error {
	_, err := x.Exec(builder.Delete(builder.NotIn("`repo_id`", builder.Select("`id`").From("`repository`"))).From("`repo_topic`"))
	if err != nil {
		return err
	}

	_, err = x.Exec(builder.Update(
		builder.Eq{
			"`repo_count`": builder.Select("count(*)").From("`repo_topic`").Where(builder.Eq{
				"`repo_topic`.`topic_id`": builder.Expr("`topic`.`id`"),
			}),
		}).From("`topic`").Where(builder.Eq{"'1'": "1"}))
	return err
}
