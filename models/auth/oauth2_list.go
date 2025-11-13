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

package auth

import (
	"github.com/kumose/kmup/models/db"

	"xorm.io/builder"
)

type FindOAuth2ApplicationsOptions struct {
	db.ListOptions
	// OwnerID is the user id or org id of the owner of the application
	OwnerID int64
	// find global applications, if true, then OwnerID will be igonred
	IsGlobal bool
}

func (opts FindOAuth2ApplicationsOptions) ToConds() builder.Cond {
	conds := builder.NewCond()
	if opts.IsGlobal {
		conds = conds.And(builder.Eq{"uid": 0})
	} else if opts.OwnerID != 0 {
		conds = conds.And(builder.Eq{"uid": opts.OwnerID})
	}
	return conds
}

func (opts FindOAuth2ApplicationsOptions) ToOrders() string {
	return "id DESC"
}
