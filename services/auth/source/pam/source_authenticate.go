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

package pam

import (
	"context"
	"fmt"
	"strings"

	"github.com/kumose/kmup/models/auth"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/auth/pam"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"

	"github.com/google/uuid"
)

// Authenticate queries if login/password is valid against the PAM,
// and create a local user if success when enabled.
func (source *Source) Authenticate(ctx context.Context, user *user_model.User, userName, password string) (*user_model.User, error) {
	pamLogin, err := pam.Auth(source.ServiceName, userName, password)
	if err != nil {
		if strings.Contains(err.Error(), "Authentication failure") {
			return nil, user_model.ErrUserNotExist{Name: userName}
		}
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	// Allow PAM sources with `@` in their name, like from Active Directory
	username := pamLogin
	email := pamLogin
	idx := strings.Index(pamLogin, "@")
	if idx > -1 {
		username = pamLogin[:idx]
	}
	if user_model.ValidateEmail(email) != nil {
		if source.EmailDomain != "" {
			email = fmt.Sprintf("%s@%s", username, source.EmailDomain)
		} else {
			email = fmt.Sprintf("%s@%s", username, setting.Service.NoReplyAddress)
		}
		if user_model.ValidateEmail(email) != nil {
			email = uuid.New().String() + "@localhost"
		}
	}

	user = &user_model.User{
		LowerName:   strings.ToLower(username),
		Name:        username,
		Email:       email,
		Passwd:      password,
		LoginType:   auth.PAM,
		LoginSource: source.AuthSource.ID,
		LoginName:   userName, // This is what the user typed in
	}
	overwriteDefault := &user_model.CreateUserOverwriteOptions{
		IsActive: optional.Some(true),
	}

	if err := user_model.CreateUser(ctx, user, &user_model.Meta{}, overwriteDefault); err != nil {
		return user, err
	}

	return user, nil
}
