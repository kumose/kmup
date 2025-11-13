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

package db

import (
	"context"
	"fmt"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

// ErrUserPasswordNotSet represents a "ErrUserPasswordNotSet" kind of error.
type ErrUserPasswordNotSet struct {
	UID  int64
	Name string
}

func (err ErrUserPasswordNotSet) Error() string {
	return fmt.Sprintf("user's password isn't set [uid: %d, name: %s]", err.UID, err.Name)
}

// Unwrap unwraps this error as a ErrInvalidArgument error
func (err ErrUserPasswordNotSet) Unwrap() error {
	return util.ErrInvalidArgument
}

// ErrUserPasswordInvalid represents a "ErrUserPasswordInvalid" kind of error.
type ErrUserPasswordInvalid struct {
	UID  int64
	Name string
}

func (err ErrUserPasswordInvalid) Error() string {
	return fmt.Sprintf("user's password is invalid [uid: %d, name: %s]", err.UID, err.Name)
}

// Unwrap unwraps this error as a ErrInvalidArgument error
func (err ErrUserPasswordInvalid) Unwrap() error {
	return util.ErrInvalidArgument
}

// Authenticate authenticates the provided user against the DB
func Authenticate(ctx context.Context, user *user_model.User, login, password string) (*user_model.User, error) {
	if user == nil {
		return nil, user_model.ErrUserNotExist{Name: login}
	}

	if !user.IsPasswordSet() {
		return nil, ErrUserPasswordNotSet{UID: user.ID, Name: user.Name}
	} else if !user.ValidatePassword(password) {
		return nil, ErrUserPasswordInvalid{UID: user.ID, Name: user.Name}
	}

	// Update password hash if server password hash algorithm have changed
	// Or update the password when the salt length doesn't match the current
	// recommended salt length, this in order to migrate user's salts to a more secure salt.
	if user.PasswdHashAlgo != setting.PasswordHashAlgo || len(user.Salt) != user_model.SaltByteLength*2 {
		if err := user.SetPassword(password); err != nil {
			return nil, err
		}
		if err := user_model.UpdateUserCols(ctx, user, "passwd", "passwd_hash_algo", "salt"); err != nil {
			return nil, err
		}
	}

	// WARN: DON'T check user.IsActive, that will be checked on reqSign so that
	// user could be hinted to resend confirm email.
	if user.ProhibitLogin {
		return nil, user_model.ErrUserProhibitLogin{
			UID:  user.ID,
			Name: user.Name,
		}
	}

	// attempting to login as a non-user account
	if user.Type != user_model.UserTypeIndividual {
		return nil, user_model.ErrUserProhibitLogin{
			UID:  user.ID,
			Name: user.Name,
		}
	}

	return user, nil
}
