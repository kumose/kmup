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
	"github.com/kumose/kmup/modules/setting"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/yandex"
)

// SimpleProviderNewFn create goth.Providers without custom url features
type SimpleProviderNewFn func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider

// SimpleProvider is a GothProvider which does not have custom url features
type SimpleProvider struct {
	BaseProvider
	scopes []string
	newFn  SimpleProviderNewFn
}

// CreateGothProvider creates a GothProvider from this Provider
func (c *SimpleProvider) CreateGothProvider(providerName, callbackURL string, source *Source) (goth.Provider, error) {
	scopes := make([]string, len(c.scopes)+len(source.Scopes))
	copy(scopes, c.scopes)
	copy(scopes[len(c.scopes):], source.Scopes)
	return c.newFn(source.ClientID, source.ClientSecret, callbackURL, scopes...), nil
}

// NewSimpleProvider is a constructor function for simple providers
func NewSimpleProvider(name, displayName string, scopes []string, newFn SimpleProviderNewFn) *SimpleProvider {
	return &SimpleProvider{
		BaseProvider: BaseProvider{
			name:        name,
			displayName: displayName,
		},
		scopes: scopes,
		newFn:  newFn,
	}
}

var _ GothProvider = &SimpleProvider{}

func init() {
	RegisterGothProvider(
		NewSimpleProvider("bitbucket", "Bitbucket", []string{"account"},
			func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
				return bitbucket.New(clientKey, secret, callbackURL, scopes...)
			}))

	RegisterGothProvider(
		NewSimpleProvider("dropbox", "Dropbox", nil,
			func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
				return dropbox.New(clientKey, secret, callbackURL, scopes...)
			}))

	RegisterGothProvider(NewSimpleProvider("facebook", "Facebook", nil,
		func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
			return facebook.New(clientKey, secret, callbackURL, scopes...)
		}))

	// named gplus due to legacy gplus -> google migration (Google killed Google+). This ensures old connections still work
	RegisterGothProvider(NewSimpleProvider("gplus", "Google", []string{"email"},
		func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
			if setting.OAuth2Client.UpdateAvatar || setting.OAuth2Client.EnableAutoRegistration {
				scopes = append(scopes, "profile")
			}
			return google.New(clientKey, secret, callbackURL, scopes...)
		}))

	RegisterGothProvider(NewSimpleProvider("twitter", "Twitter", nil,
		func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
			return twitter.New(clientKey, secret, callbackURL)
		}))

	RegisterGothProvider(NewSimpleProvider("discord", "Discord", []string{discord.ScopeIdentify, discord.ScopeEmail},
		func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
			return discord.New(clientKey, secret, callbackURL, scopes...)
		}))

	// See https://tech.yandex.com/passport/doc/dg/reference/response-docpage/
	RegisterGothProvider(NewSimpleProvider("yandex", "Yandex", []string{"login:email", "login:info", "login:avatar"},
		func(clientKey, secret, callbackURL string, scopes ...string) goth.Provider {
			return yandex.New(clientKey, secret, callbackURL, scopes...)
		}))

	RegisterGothProvider(NewSimpleProvider(
		"azuread", "Azure AD", nil,
		func(clientID, secret, callbackURL string, scopes ...string) goth.Provider {
			return azuread.New(clientID, secret, callbackURL, nil, scopes...)
		},
	))

	RegisterGothProvider(NewSimpleProvider(
		"microsoftonline", "Microsoft Online", nil,
		func(clientID, secret, callbackURL string, scopes ...string) goth.Provider {
			return microsoftonline.New(clientID, secret, callbackURL, scopes...)
		},
	))
}
