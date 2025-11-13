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

	"github.com/kumose/kmup/models/auth"
	user_model "github.com/kumose/kmup/models/user"
)

// Source is a password authentication service
type Source struct {
	auth.ConfigBase `json:"-"`
}

// FromDB fills up an OAuth2Config from serialized format.
func (source *Source) FromDB(bs []byte) error {
	return nil
}

// ToDB exports the config to a byte slice to be saved into database (this method is just dummy and does nothing for DB source)
func (source *Source) ToDB() ([]byte, error) {
	return nil, nil
}

// Authenticate queries if login/password is valid against the PAM,
// and create a local user if success when enabled.
func (source *Source) Authenticate(ctx context.Context, user *user_model.User, login, password string) (*user_model.User, error) {
	return Authenticate(ctx, user, login, password)
}

func init() {
	auth.RegisterTypeConfig(auth.NoType, &Source{})
	auth.RegisterTypeConfig(auth.Plain, &Source{})
}
