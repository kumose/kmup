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

package oauth2

import (
	"fmt"
)

// RegisterSource causes an OAuth2 configuration to be registered
func (source *Source) RegisterSource() error {
	err := RegisterProviderWithGothic(source.AuthSource.Name, source)
	return wrapOpenIDConnectInitializeError(err, source.AuthSource.Name, source)
}

// UnregisterSource causes an OAuth2 configuration to be unregistered
func (source *Source) UnregisterSource() error {
	RemoveProviderFromGothic(source.AuthSource.Name)
	return nil
}

// ErrOpenIDConnectInitialize represents a "OpenIDConnectInitialize" kind of error.
type ErrOpenIDConnectInitialize struct {
	OpenIDConnectAutoDiscoveryURL string
	ProviderName                  string
	Cause                         error
}

// IsErrOpenIDConnectInitialize checks if an error is a ExternalLoginUserAlreadyExist.
func IsErrOpenIDConnectInitialize(err error) bool {
	_, ok := err.(ErrOpenIDConnectInitialize)
	return ok
}

func (err ErrOpenIDConnectInitialize) Error() string {
	return fmt.Sprintf("Failed to initialize OpenID Connect Provider with name '%s' with url '%s': %v", err.ProviderName, err.OpenIDConnectAutoDiscoveryURL, err.Cause)
}

func (err ErrOpenIDConnectInitialize) Unwrap() error {
	return err.Cause
}

// wrapOpenIDConnectInitializeError is used to wrap the error but this cannot be done in modules/auth/oauth2
// inside oauth2: import cycle not allowed models -> modules/auth/oauth2 -> models
func wrapOpenIDConnectInitializeError(err error, providerName string, source *Source) error {
	if err != nil && source.Provider == "openidConnect" {
		err = ErrOpenIDConnectInitialize{ProviderName: providerName, OpenIDConnectAutoDiscoveryURL: source.OpenIDConnectAutoDiscoveryURL, Cause: err}
	}
	return err
}
