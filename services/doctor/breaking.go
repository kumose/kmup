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

package doctor

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"

	"xorm.io/builder"
)

func iterateUserAccounts(ctx context.Context, each func(*user.User) error) error {
	err := db.Iterate(
		ctx,
		builder.Gt{"id": 0},
		func(ctx context.Context, bean *user.User) error {
			return each(bean)
		},
	)
	return err
}

// Since 1.16.4 new restrictions has been set on email addresses. However users with invalid email
// addresses would be currently facing a error due to their invalid email address.
// Ref: https://github.com/go-kmup/kmup/pull/19085 & https://github.com/go-kmup/kmup/pull/17688
func checkUserEmail(ctx context.Context, logger log.Logger, _ bool) error {
	// We could use quirky SQL to get all users that start without a [a-zA-Z0-9], but that would mean
	// DB provider-specific SQL and only works _now_. So instead we iterate through all user accounts
	// and use the user.ValidateEmail function to be future-proof.
	var invalidUserCount int64
	if err := iterateUserAccounts(ctx, func(u *user.User) error {
		// Only check for users, skip
		if u.Type != user.UserTypeIndividual {
			return nil
		}

		if err := user.ValidateEmail(u.Email); err != nil {
			invalidUserCount++
			logger.Warn("User[id=%d name=%q] have not a valid e-mail: %v", u.ID, u.Name, err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("iterateUserAccounts: %w", err)
	}

	if invalidUserCount == 0 {
		logger.Info("All users have a valid e-mail.")
	} else {
		logger.Warn("%d user(s) have a non-valid e-mail.", invalidUserCount)
	}
	return nil
}

// From time to time Kmup makes changes to the reserved usernames and which symbols
// are allowed for various reasons. This check helps with detecting users that, according
// to our reserved names, don't have a valid username.
func checkUserName(ctx context.Context, logger log.Logger, _ bool) error {
	var invalidUserCount int64
	if err := iterateUserAccounts(ctx, func(u *user.User) error {
		if err := user.IsUsableUsername(u.Name); err != nil {
			invalidUserCount++
			logger.Warn("User[id=%d] does not have a valid username: %v", u.ID, err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("iterateUserAccounts: %w", err)
	}

	if invalidUserCount == 0 {
		logger.Info("All users have a valid username.")
	} else {
		logger.Warn("%d user(s) have a non-valid username.", invalidUserCount)
	}
	return nil
}

func init() {
	Register(&Check{
		Title:     "Check if users has an valid email address",
		Name:      "check-user-email",
		IsDefault: false,
		Run:       checkUserEmail,
		Priority:  9,
	})
	Register(&Check{
		Title:     "Check if users have a valid username",
		Name:      "check-user-names",
		IsDefault: false,
		Run:       checkUserName,
		Priority:  9,
	})
}
