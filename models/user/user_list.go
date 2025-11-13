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

	"github.com/kumose/kmup/models/db"
)

func GetUsersMapByIDs(ctx context.Context, userIDs []int64) (map[int64]*User, error) {
	userMaps := make(map[int64]*User, len(userIDs))
	if len(userIDs) == 0 {
		return userMaps, nil
	}

	left := len(userIDs)
	for left > 0 {
		limit := min(left, db.DefaultMaxInSize)
		err := db.GetEngine(ctx).
			In("id", userIDs[:limit]).
			Find(&userMaps)
		if err != nil {
			return nil, err
		}
		left -= limit
		userIDs = userIDs[limit:]
	}
	return userMaps, nil
}

func GetPossibleUserFromMap(userID int64, usererMaps map[int64]*User) *User {
	switch userID {
	case GhostUserID:
		return NewGhostUser()
	case ActionsUserID:
		return NewActionsUser()
	case 0:
		return nil
	default:
		user, ok := usererMaps[userID]
		if !ok {
			return NewGhostUser()
		}
		return user
	}
}
