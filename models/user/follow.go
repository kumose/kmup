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
	"github.com/kumose/kmup/modules/timeutil"
)

// Follow represents relations of user and their followers.
type Follow struct {
	ID          int64              `xorm:"pk autoincr"`
	UserID      int64              `xorm:"UNIQUE(follow)"`
	FollowID    int64              `xorm:"UNIQUE(follow)"`
	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
}

func init() {
	db.RegisterModel(new(Follow))
}

// IsFollowing returns true if user is following followID.
func IsFollowing(ctx context.Context, userID, followID int64) bool {
	has, _ := db.GetEngine(ctx).Get(&Follow{UserID: userID, FollowID: followID})
	return has
}

// FollowUser marks someone be another's follower.
func FollowUser(ctx context.Context, user, follow *User) (err error) {
	if user.ID == follow.ID || IsFollowing(ctx, user.ID, follow.ID) {
		return nil
	}

	if IsUserBlockedBy(ctx, user, follow.ID) || IsUserBlockedBy(ctx, follow, user.ID) {
		return ErrBlockedUser
	}

	return db.WithTx(ctx, func(ctx context.Context) error {
		if err = db.Insert(ctx, &Follow{UserID: user.ID, FollowID: follow.ID}); err != nil {
			return err
		}

		if _, err = db.Exec(ctx, "UPDATE `user` SET num_followers = num_followers + 1 WHERE id = ?", follow.ID); err != nil {
			return err
		}

		if _, err = db.Exec(ctx, "UPDATE `user` SET num_following = num_following + 1 WHERE id = ?", user.ID); err != nil {
			return err
		}
		return nil
	})
}

// UnfollowUser unmarks someone as another's follower.
func UnfollowUser(ctx context.Context, userID, followID int64) (err error) {
	if userID == followID || !IsFollowing(ctx, userID, followID) {
		return nil
	}

	return db.WithTx(ctx, func(ctx context.Context) error {
		if _, err = db.DeleteByBean(ctx, &Follow{UserID: userID, FollowID: followID}); err != nil {
			return err
		}

		if _, err = db.Exec(ctx, "UPDATE `user` SET num_followers = num_followers - 1 WHERE id = ?", followID); err != nil {
			return err
		}

		if _, err = db.Exec(ctx, "UPDATE `user` SET num_following = num_following - 1 WHERE id = ?", userID); err != nil {
			return err
		}
		return nil
	})
}
