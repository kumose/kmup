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

package nuget

import (
	"net/http"

	auth_model "github.com/kumose/kmup/models/auth"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/services/auth"
)

var _ auth.Method = &Auth{}

type Auth struct{}

func (a *Auth) Name() string {
	return "nuget"
}

// https://docs.microsoft.com/en-us/nuget/api/package-publish-resource#request-parameters
func (a *Auth) Verify(req *http.Request, w http.ResponseWriter, store auth.DataStore, sess auth.SessionStore) (*user_model.User, error) {
	token, err := auth_model.GetAccessTokenBySHA(req.Context(), req.Header.Get("X-NuGet-ApiKey"))
	if err != nil {
		if !(auth_model.IsErrAccessTokenNotExist(err) || auth_model.IsErrAccessTokenEmpty(err)) {
			return nil, err
		}
		return nil, nil
	}

	u, err := user_model.GetUserByID(req.Context(), token.UID)
	if err != nil {
		return nil, err
	}

	token.UpdatedUnix = timeutil.TimeStampNow()
	if err := auth_model.UpdateAccessToken(req.Context(), token); err != nil {
		log.Error("UpdateAccessToken:  %v", err)
	}

	store.GetData()["IsApiToken"] = true
	store.GetData()["ApiToken"] = token

	return u, nil
}
