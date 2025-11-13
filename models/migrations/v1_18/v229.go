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

package v1_18

import (
	"fmt"

	"github.com/kumose/kmup/models/issues"

	"xorm.io/builder"
	"xorm.io/xorm"
)

func UpdateOpenMilestoneCounts(x *xorm.Engine) error {
	var openMilestoneIDs []int64
	err := x.Table("milestone").Select("id").Where(builder.Neq{"is_closed": 1}).Find(&openMilestoneIDs)
	if err != nil {
		return fmt.Errorf("error selecting open milestone IDs: %w", err)
	}

	for _, id := range openMilestoneIDs {
		_, err := x.ID(id).
			Cols("num_issues", "num_closed_issues").
			SetExpr("num_issues", builder.Select("count(*)").From("issue").Where(
				builder.Eq{"milestone_id": id},
			)).
			SetExpr("num_closed_issues", builder.Select("count(*)").From("issue").Where(
				builder.Eq{
					"milestone_id": id,
					"is_closed":    true,
				},
			)).
			Update(&issues.Milestone{})
		if err != nil {
			return fmt.Errorf("error updating issue counts in milestone %d: %w", id, err)
		}
		_, err = x.Exec("UPDATE `milestone` SET completeness=100*num_closed_issues/(CASE WHEN num_issues > 0 THEN num_issues ELSE 1 END) WHERE id=?",
			id,
		)
		if err != nil {
			return fmt.Errorf("error setting completeness on milestone %d: %w", id, err)
		}
	}

	return nil
}
