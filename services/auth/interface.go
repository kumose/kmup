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
	"net/http"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/reqctx"
	"github.com/kumose/kmup/modules/session"
)

type DataStore = reqctx.ContextDataProvider

// SessionStore represents a session store
type SessionStore session.Store

// Method represents an authentication method (plugin) for HTTP requests.
type Method interface {
	// Verify tries to verify the authentication data contained in the request.
	// If verification succeeds, it returns either an existing user object (with id > 0)
	// or a new user object (with id = 0) populated with the information that was found
	// in the authentication data (username or email).
	// Second argument returns err if verification fails, otherwise
	// First return argument returns nil if no matched verification condition
	Verify(http *http.Request, w http.ResponseWriter, store DataStore, sess SessionStore) (*user_model.User, error)

	Name() string
}

// PasswordAuthenticator represents a source of authentication
type PasswordAuthenticator interface {
	Authenticate(ctx context.Context, user *user_model.User, login, password string) (*user_model.User, error)
}

// SynchronizableSource represents a source that can synchronize users
type SynchronizableSource interface {
	Sync(ctx context.Context, updateExisting bool) error
}
