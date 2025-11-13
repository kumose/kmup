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
	"slices"
	"strconv"

	"github.com/kumose/kmup/models/user"
)

func MakeSelfOnTop(doer *user.User, users []*user.User) []*user.User {
	if doer != nil {
		idx := slices.IndexFunc(users, func(u *user.User) bool {
			return u.ID == doer.ID
		})
		if idx > 0 {
			newUsers := make([]*user.User, len(users))
			newUsers[0] = users[idx]
			copy(newUsers[1:], users[:idx])
			copy(newUsers[idx+1:], users[idx+1:])
			return newUsers
		}
	}
	return users
}

// GetFilterUserIDByName tries to get the user ID from the given username.
// Before, the "issue filter" passes user ID to query the list, but in many cases, it's impossible to pre-fetch the full user list.
// So it's better to make it work like GitHub: users could input username directly.
// Since it only converts the username to ID directly and is only used internally (to search issues), so no permission check is needed.
// Return values:
// * "": no filter
// * "{the-id}": match the id
// * "(none)": match no issue (due to the user doesn't exist)
func GetFilterUserIDByName(ctx context.Context, name string) string {
	if name == "" {
		return ""
	}
	u, err := user.GetUserByName(ctx, name)
	if err != nil {
		if id, err := strconv.ParseInt(name, 10, 64); err == nil {
			return strconv.FormatInt(id, 10)
		}
		// The "(none)" is for internal usage only: when doer tries to search non-existing user, use "(none)" to return empty result.
		return "(none)"
	}
	return strconv.FormatInt(u.ID, 10)
}
