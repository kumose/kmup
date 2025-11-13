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
	"fmt"
	"strings"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/util"
)

// ErrUserRedirectNotExist represents a "UserRedirectNotExist" kind of error.
type ErrUserRedirectNotExist struct {
	Name string
}

// IsErrUserRedirectNotExist check if an error is an ErrUserRedirectNotExist.
func IsErrUserRedirectNotExist(err error) bool {
	_, ok := err.(ErrUserRedirectNotExist)
	return ok
}

func (err ErrUserRedirectNotExist) Error() string {
	return fmt.Sprintf("user redirect does not exist [name: %s]", err.Name)
}

func (err ErrUserRedirectNotExist) Unwrap() error {
	return util.ErrNotExist
}

// Redirect represents that a user name should be redirected to another
type Redirect struct {
	ID             int64  `xorm:"pk autoincr"`
	LowerName      string `xorm:"UNIQUE(s) INDEX NOT NULL"`
	RedirectUserID int64  // userID to redirect to
}

// TableName provides the real table name
func (Redirect) TableName() string {
	return "user_redirect"
}

func init() {
	db.RegisterModel(new(Redirect))
}

// LookupUserRedirect look up userID if a user has a redirect name
func LookupUserRedirect(ctx context.Context, userName string) (int64, error) {
	userName = strings.ToLower(userName)
	redirect := &Redirect{LowerName: userName}
	if has, err := db.GetEngine(ctx).Get(redirect); err != nil {
		return 0, err
	} else if !has {
		return 0, ErrUserRedirectNotExist{Name: userName}
	}
	return redirect.RedirectUserID, nil
}

// NewUserRedirect create a new user redirect
func NewUserRedirect(ctx context.Context, ID int64, oldUserName, newUserName string) error {
	oldUserName = strings.ToLower(oldUserName)
	newUserName = strings.ToLower(newUserName)

	if err := DeleteUserRedirect(ctx, oldUserName); err != nil {
		return err
	}

	if err := DeleteUserRedirect(ctx, newUserName); err != nil {
		return err
	}

	return db.Insert(ctx, &Redirect{
		LowerName:      oldUserName,
		RedirectUserID: ID,
	})
}

// DeleteUserRedirect delete any redirect from the specified user name to
// anything else
func DeleteUserRedirect(ctx context.Context, userName string) error {
	userName = strings.ToLower(userName)
	_, err := db.GetEngine(ctx).Delete(&Redirect{LowerName: userName})
	return err
}
