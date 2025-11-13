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

package user

import (
	"context"
	"strings"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/builder"
)

func SetMustChangePassword(ctx context.Context, all, mustChangePassword bool, include, exclude []string) (int64, error) {
	sliceTrimSpaceDropEmpty := func(input []string) []string {
		output := make([]string, 0, len(input))
		for _, in := range input {
			in = strings.ToLower(strings.TrimSpace(in))
			if in == "" {
				continue
			}
			output = append(output, in)
		}
		return output
	}

	var cond builder.Cond

	// Only include the users where something changes to get an accurate count
	cond = builder.Neq{"must_change_password": mustChangePassword}

	if !all {
		include = sliceTrimSpaceDropEmpty(include)
		if len(include) == 0 {
			return 0, util.ErrorWrap(util.ErrInvalidArgument, "no users to include provided")
		}

		cond = cond.And(builder.In("lower_name", include))
	}

	exclude = sliceTrimSpaceDropEmpty(exclude)
	if len(exclude) > 0 {
		cond = cond.And(builder.NotIn("lower_name", exclude))
	}

	return db.GetEngine(ctx).Where(cond).MustCols("must_change_password").Update(&User{MustChangePassword: mustChangePassword})
}
