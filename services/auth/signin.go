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

package auth

import (
	"context"
	"strings"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/services/auth/source/oauth2"
	"github.com/kumose/kmup/services/auth/source/smtp"

	_ "github.com/kumose/kmup/services/auth/source/db"   // register the sources (and below)
	_ "github.com/kumose/kmup/services/auth/source/ldap" // register the ldap source
	_ "github.com/kumose/kmup/services/auth/source/pam"  // register the pam source
	_ "github.com/kumose/kmup/services/auth/source/sspi" // register the sspi source
)

// UserSignIn validates user name and password.
func UserSignIn(ctx context.Context, username, password string) (*user_model.User, *auth.Source, error) {
	var user *user_model.User
	isEmail := false
	if strings.Contains(username, "@") {
		isEmail = true
		emailAddress := user_model.EmailAddress{LowerEmail: strings.ToLower(strings.TrimSpace(username))}
		// check same email
		has, err := db.GetEngine(ctx).Get(&emailAddress)
		if err != nil {
			return nil, nil, err
		}
		if has {
			if !emailAddress.IsActivated {
				return nil, nil, user_model.ErrEmailAddressNotExist{
					Email: username,
				}
			}
			user = &user_model.User{ID: emailAddress.UID}
		}
	} else {
		trimmedUsername := strings.TrimSpace(username)
		if len(trimmedUsername) == 0 {
			return nil, nil, user_model.ErrUserNotExist{Name: username}
		}

		user = &user_model.User{LowerName: strings.ToLower(trimmedUsername)}
	}

	if user != nil {
		hasUser, err := user_model.GetUser(ctx, user)
		if err != nil {
			return nil, nil, err
		}

		if hasUser {
			source, err := auth.GetSourceByID(ctx, user.LoginSource)
			if err != nil {
				return nil, nil, err
			}

			if !source.IsActive {
				return nil, nil, oauth2.ErrAuthSourceNotActivated
			}

			authenticator, ok := source.Cfg.(PasswordAuthenticator)
			if !ok {
				return nil, nil, smtp.ErrUnsupportedLoginType
			}

			user, err := authenticator.Authenticate(ctx, user, user.LoginName, password)
			if err != nil {
				return nil, nil, err
			}

			// WARN: DON'T check user.IsActive, that will be checked on reqSign so that
			// user could be hint to resend confirm email.
			if user.ProhibitLogin {
				return nil, nil, user_model.ErrUserProhibitLogin{UID: user.ID, Name: user.Name}
			}

			return user, source, nil
		}
	}

	sources, err := db.Find[auth.Source](ctx, auth.FindSourcesOptions{
		IsActive: optional.Some(true),
	})
	if err != nil {
		return nil, nil, err
	}

	for _, source := range sources {
		if !source.IsActive {
			// don't try to authenticate non-active sources
			continue
		}

		authenticator, ok := source.Cfg.(PasswordAuthenticator)
		if !ok {
			continue
		}

		authUser, err := authenticator.Authenticate(ctx, nil, username, password)

		if err == nil {
			if !authUser.ProhibitLogin {
				return authUser, source, nil
			}
			err = user_model.ErrUserProhibitLogin{UID: authUser.ID, Name: authUser.Name}
		}

		if user_model.IsErrUserNotExist(err) {
			log.Debug("Failed to login '%s' via '%s': %v", username, source.Name, err)
		} else {
			log.Warn("Failed to login '%s' via '%s': %v", username, source.Name, err)
		}
	}

	if isEmail {
		return nil, nil, user_model.ErrEmailAddressNotExist{Email: username}
	}

	return nil, nil, user_model.ErrUserNotExist{Name: username}
}
