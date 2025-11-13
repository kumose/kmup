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

package password

import (
	"context"
	"errors"
	"fmt"

	"github.com/kumose/kmup/modules/auth/password/pwn"
	"github.com/kumose/kmup/modules/setting"
)

var ErrIsPwned = errors.New("password has been pwned")

type ErrIsPwnedRequest struct {
	err error
}

func IsErrIsPwnedRequest(err error) bool {
	_, ok := err.(ErrIsPwnedRequest)
	return ok
}

func (err ErrIsPwnedRequest) Error() string {
	return fmt.Sprintf("using Have-I-Been-Pwned service failed: %v", err.err)
}

func (err ErrIsPwnedRequest) Unwrap() error {
	return err.err
}

// IsPwned checks whether a password has been pwned
// If a password has not been pwned, no error is returned.
func IsPwned(ctx context.Context, password string) error {
	if !setting.PasswordCheckPwn {
		return nil
	}

	client := pwn.New(pwn.WithContext(ctx))
	count, err := client.CheckPassword(password, true)
	if err != nil {
		return ErrIsPwnedRequest{err}
	}

	if count > 0 {
		return ErrIsPwned
	}

	return nil
}
