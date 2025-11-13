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

package conan

import (
	"net/http"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/auth"
	"github.com/kumose/kmup/services/packages"
)

var _ auth.Method = &Auth{}

type Auth struct{}

func (a *Auth) Name() string {
	return "conan"
}

// Verify extracts the user from the Bearer token
func (a *Auth) Verify(req *http.Request, w http.ResponseWriter, store auth.DataStore, sess auth.SessionStore) (*user_model.User, error) {
	packageMeta, err := packages.ParseAuthorizationRequest(req)
	if err != nil {
		log.Trace("ParseAuthorizationToken: %v", err)
		return nil, err
	}

	if packageMeta == nil || packageMeta.UserID == 0 {
		return nil, nil
	}

	u, err := user_model.GetUserByID(req.Context(), packageMeta.UserID)
	if err != nil {
		return nil, err
	}
	if packageMeta.Scope != "" {
		store.GetData()["IsApiToken"] = true
		store.GetData()["ApiTokenScope"] = packageMeta.Scope
	}

	return u, nil
}
