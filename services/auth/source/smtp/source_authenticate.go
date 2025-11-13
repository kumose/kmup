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

package smtp

import (
	"context"
	"errors"
	"net/smtp"
	"net/textproto"
	"strings"

	auth_model "github.com/kumose/kmup/models/auth"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/util"
)

// Authenticate queries if the provided login/password is authenticates against the SMTP server
// Users will be autoregistered as required
func (source *Source) Authenticate(ctx context.Context, user *user_model.User, userName, password string) (*user_model.User, error) {
	// Verify allowed domains.
	if len(source.AllowedDomains) > 0 {
		idx := strings.Index(userName, "@")
		if idx == -1 {
			return nil, user_model.ErrUserNotExist{Name: userName}
		} else if !util.SliceContainsString(strings.Split(source.AllowedDomains, ","), userName[idx+1:], true) {
			return nil, user_model.ErrUserNotExist{Name: userName}
		}
	}

	var auth smtp.Auth
	switch source.Auth {
	case PlainAuthentication:
		auth = smtp.PlainAuth("", userName, password, source.Host)
	case LoginAuthentication:
		auth = &loginAuthenticator{userName, password}
	case CRAMMD5Authentication:
		auth = smtp.CRAMMD5Auth(userName, password)
	default:
		return nil, errors.New("unsupported SMTP auth type")
	}

	if err := Authenticate(auth, source); err != nil {
		// Check standard error format first,
		// then fallback to worse case.
		tperr, ok := err.(*textproto.Error)
		if (ok && tperr.Code == 535) ||
			strings.Contains(err.Error(), "Username and Password not accepted") {
			return nil, user_model.ErrUserNotExist{Name: userName}
		}
		if (ok && tperr.Code == 534) ||
			strings.Contains(err.Error(), "Application-specific password required") {
			return nil, user_model.ErrUserNotExist{Name: userName}
		}
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	username := userName
	idx := strings.Index(userName, "@")
	if idx > -1 {
		username = userName[:idx]
	}

	user = &user_model.User{
		LowerName:   strings.ToLower(username),
		Name:        strings.ToLower(username),
		Email:       userName,
		Passwd:      password,
		LoginType:   auth_model.SMTP,
		LoginSource: source.AuthSource.ID,
		LoginName:   userName,
	}
	overwriteDefault := &user_model.CreateUserOverwriteOptions{
		IsActive: optional.Some(true),
	}

	if err := user_model.CreateUser(ctx, user, &user_model.Meta{}, overwriteDefault); err != nil {
		return user, err
	}

	return user, nil
}
