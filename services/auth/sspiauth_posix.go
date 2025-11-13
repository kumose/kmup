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

//go:build !windows

package auth

import (
	"errors"
	"net/http"
)

type SSPIUserInfo struct {
	Username string   // Name of user, usually in the form DOMAIN\User
	Groups   []string // The global groups the user is a member of
}

type sspiAuthMock struct{}

func (s sspiAuthMock) AppendAuthenticateHeader(w http.ResponseWriter, data string) {
}

func (s sspiAuthMock) Authenticate(r *http.Request, w http.ResponseWriter) (userInfo *SSPIUserInfo, outToken string, err error) {
	return nil, "", errors.New("not implemented")
}

func sspiAuthInit() error {
	sspiAuth = &sspiAuthMock{} // TODO: we can mock the SSPI auth in tests
	return nil
}
