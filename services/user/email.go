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
	"errors"
	"strings"

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

// AdminAddOrSetPrimaryEmailAddress is used by admins to add or set a user's primary email address
func AdminAddOrSetPrimaryEmailAddress(ctx context.Context, u *user_model.User, emailStr string) error {
	if strings.EqualFold(u.Email, emailStr) {
		return nil
	}

	if err := user_model.ValidateEmailForAdmin(emailStr); err != nil {
		return err
	}

	// Check if address exists already
	email, err := user_model.GetEmailAddressByEmail(ctx, emailStr)
	if err != nil && !errors.Is(err, util.ErrNotExist) {
		return err
	}
	if email != nil && email.UID != u.ID {
		return user_model.ErrEmailAlreadyUsed{Email: emailStr}
	}

	// Update old primary address
	primary, err := user_model.GetPrimaryEmailAddressOfUser(ctx, u.ID)
	if err != nil {
		return err
	}

	primary.IsPrimary = false
	if err := user_model.UpdateEmailAddress(ctx, primary); err != nil {
		return err
	}

	// Insert new or update existing address
	if email != nil {
		email.IsPrimary = true
		email.IsActivated = true
		if err := user_model.UpdateEmailAddress(ctx, email); err != nil {
			return err
		}
	} else {
		email = &user_model.EmailAddress{
			UID:         u.ID,
			Email:       emailStr,
			IsActivated: true,
			IsPrimary:   true,
		}
		if _, err := user_model.InsertEmailAddress(ctx, email); err != nil {
			return err
		}
	}

	u.Email = emailStr

	return user_model.UpdateUserCols(ctx, u, "email")
}

func ReplacePrimaryEmailAddress(ctx context.Context, u *user_model.User, emailStr string) error {
	if strings.EqualFold(u.Email, emailStr) {
		return nil
	}

	if err := user_model.ValidateEmail(emailStr); err != nil {
		return err
	}

	if !u.IsOrganization() {
		// Check if address exists already
		email, err := user_model.GetEmailAddressByEmail(ctx, emailStr)
		if err != nil && !errors.Is(err, util.ErrNotExist) {
			return err
		}
		if email != nil {
			if email.IsPrimary && email.UID == u.ID {
				return nil
			}
			return user_model.ErrEmailAlreadyUsed{Email: emailStr}
		}

		// Remove old primary address
		primary, err := user_model.GetPrimaryEmailAddressOfUser(ctx, u.ID)
		if err != nil {
			return err
		}
		if _, err := db.DeleteByID[user_model.EmailAddress](ctx, primary.ID); err != nil {
			return err
		}

		// Insert new primary address
		email = &user_model.EmailAddress{
			UID:         u.ID,
			Email:       emailStr,
			IsActivated: true,
			IsPrimary:   true,
		}
		if _, err := user_model.InsertEmailAddress(ctx, email); err != nil {
			return err
		}
	}

	u.Email = emailStr

	return user_model.UpdateUserCols(ctx, u, "email")
}

func AddEmailAddresses(ctx context.Context, u *user_model.User, emails []string) error {
	for _, emailStr := range emails {
		if err := user_model.ValidateEmail(emailStr); err != nil {
			return err
		}

		// Check if address exists already
		email, err := user_model.GetEmailAddressByEmail(ctx, emailStr)
		if err != nil && !errors.Is(err, util.ErrNotExist) {
			return err
		}
		if email != nil {
			return user_model.ErrEmailAlreadyUsed{Email: emailStr}
		}

		// Insert new address
		email = &user_model.EmailAddress{
			UID:         u.ID,
			Email:       emailStr,
			IsActivated: !setting.Service.RegisterEmailConfirm,
			IsPrimary:   false,
		}
		if _, err := user_model.InsertEmailAddress(ctx, email); err != nil {
			return err
		}
	}

	return nil
}

func DeleteEmailAddresses(ctx context.Context, u *user_model.User, emails []string) error {
	for _, emailStr := range emails {
		// Check if address exists
		email, err := user_model.GetEmailAddressOfUser(ctx, emailStr, u.ID)
		if err != nil {
			return err
		}
		if email.IsPrimary {
			return user_model.ErrPrimaryEmailCannotDelete{Email: emailStr}
		}

		// Remove address
		if _, err := db.DeleteByID[user_model.EmailAddress](ctx, email.ID); err != nil {
			return err
		}
	}

	return nil
}
