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

package git

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/container"
)

// LFSLockList is a list of LFSLock
type LFSLockList []*LFSLock

// LoadAttributes loads the attributes for the given locks
func (locks LFSLockList) LoadAttributes(ctx context.Context) error {
	if len(locks) == 0 {
		return nil
	}

	if err := locks.LoadOwner(ctx); err != nil {
		return fmt.Errorf("load owner: %w", err)
	}

	return nil
}

// LoadOwner loads the owner of the locks
func (locks LFSLockList) LoadOwner(ctx context.Context) error {
	if len(locks) == 0 {
		return nil
	}

	usersIDs := container.FilterSlice(locks, func(lock *LFSLock) (int64, bool) {
		return lock.OwnerID, true
	})
	users := make(map[int64]*user_model.User, len(usersIDs))
	if err := db.GetEngine(ctx).
		In("id", usersIDs).
		Find(&users); err != nil {
		return fmt.Errorf("find users: %w", err)
	}
	for _, v := range locks {
		v.Owner = users[v.OwnerID]
		if v.Owner == nil { // not exist
			v.Owner = user_model.NewGhostUser()
		}
	}

	return nil
}
