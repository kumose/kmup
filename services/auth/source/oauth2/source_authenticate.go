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
	"context"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/services/auth/source/db"
)

// Authenticate falls back to the db authenticator
func (source *Source) Authenticate(ctx context.Context, user *user_model.User, login, password string) (*user_model.User, error) {
	return db.Authenticate(ctx, user, login, password)
}

// NB: Oauth2 does not implement LocalTwoFASkipper for password authentication
// as its password authentication drops to db authentication
