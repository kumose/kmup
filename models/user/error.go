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
	"fmt"

	"github.com/kumose/kmup/modules/util"
)

// ErrUserAlreadyExist represents a "user already exists" error.
type ErrUserAlreadyExist struct {
	Name string
}

// IsErrUserAlreadyExist checks if an error is a ErrUserAlreadyExists.
func IsErrUserAlreadyExist(err error) bool {
	_, ok := err.(ErrUserAlreadyExist)
	return ok
}

func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("user already exists [name: %s]", err.Name)
}

// Unwrap unwraps this error as a ErrExist error
func (err ErrUserAlreadyExist) Unwrap() error {
	return util.ErrAlreadyExist
}

// ErrUserNotExist represents a "UserNotExist" kind of error.
type ErrUserNotExist struct {
	UID  int64
	Name string
}

// IsErrUserNotExist checks if an error is a ErrUserNotExist.
func IsErrUserNotExist(err error) bool {
	_, ok := err.(ErrUserNotExist)
	return ok
}

func (err ErrUserNotExist) Error() string {
	return fmt.Sprintf("user does not exist [uid: %d, name: %s]", err.UID, err.Name)
}

// Unwrap unwraps this error as a ErrNotExist error
func (err ErrUserNotExist) Unwrap() error {
	return util.ErrNotExist
}

// ErrUserProhibitLogin represents a "ErrUserProhibitLogin" kind of error.
type ErrUserProhibitLogin struct {
	UID  int64
	Name string
}

// IsErrUserProhibitLogin checks if an error is a ErrUserProhibitLogin
func IsErrUserProhibitLogin(err error) bool {
	_, ok := err.(ErrUserProhibitLogin)
	return ok
}

func (err ErrUserProhibitLogin) Error() string {
	return fmt.Sprintf("user is not allowed login [uid: %d, name: %s]", err.UID, err.Name)
}

// Unwrap unwraps this error as a ErrPermission error
func (err ErrUserProhibitLogin) Unwrap() error {
	return util.ErrPermissionDenied
}

// ErrUserInactive represents a "ErrUserInactive" kind of error.
type ErrUserInactive struct {
	UID  int64
	Name string
}

// IsErrUserInactive checks if an error is a ErrUserInactive
func IsErrUserInactive(err error) bool {
	_, ok := err.(ErrUserInactive)
	return ok
}

func (err ErrUserInactive) Error() string {
	return fmt.Sprintf("user is inactive [uid: %d, name: %s]", err.UID, err.Name)
}

// Unwrap unwraps this error as a ErrPermission error
func (err ErrUserInactive) Unwrap() error {
	return util.ErrPermissionDenied
}

// ErrUserIsNotLocal represents a "ErrUserIsNotLocal" kind of error.
type ErrUserIsNotLocal struct {
	UID  int64
	Name string
}

func (err ErrUserIsNotLocal) Error() string {
	return fmt.Sprintf("user is not local type [uid: %d, name: %s]", err.UID, err.Name)
}

// IsErrUserIsNotLocal
func IsErrUserIsNotLocal(err error) bool {
	_, ok := err.(ErrUserIsNotLocal)
	return ok
}
