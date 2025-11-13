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
	"time"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

type fakeProvider struct{}

func (p *fakeProvider) Name() string {
	return "fake"
}

func (p *fakeProvider) SetName(name string) {}

func (p *fakeProvider) BeginAuth(state string) (goth.Session, error) {
	return nil, nil
}

func (p *fakeProvider) UnmarshalSession(string) (goth.Session, error) {
	return nil, nil
}

func (p *fakeProvider) FetchUser(goth.Session) (goth.User, error) {
	return goth.User{}, nil
}

func (p *fakeProvider) Debug(bool) {
}

func (p *fakeProvider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	switch refreshToken {
	case "expired":
		return nil, &oauth2.RetrieveError{
			ErrorCode: "invalid_grant",
		}
	default:
		return &oauth2.Token{
			AccessToken:  "token",
			TokenType:    "Bearer",
			RefreshToken: "refresh",
			Expiry:       time.Now().Add(time.Hour),
		}, nil
	}
}

func (p *fakeProvider) RefreshTokenAvailable() bool {
	return true
}

func init() {
	RegisterGothProvider(
		NewSimpleProvider("fake", "Fake", []string{"account"},
			func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
				return &fakeProvider{}
			}))
}
