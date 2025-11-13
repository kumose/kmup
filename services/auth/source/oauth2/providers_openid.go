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
	"html/template"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/svg"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

// OpenIDProvider is a GothProvider for OpenID
type OpenIDProvider struct{}

func (o *OpenIDProvider) SupportSSHPublicKey() bool {
	return true
}

// Name provides the technical name for this provider
func (o *OpenIDProvider) Name() string {
	return "openidConnect"
}

// DisplayName returns the friendly name for this provider
func (o *OpenIDProvider) DisplayName() string {
	return "OpenID Connect"
}

// IconHTML returns icon HTML for this provider
func (o *OpenIDProvider) IconHTML(size int) template.HTML {
	return svg.RenderHTML("kmup-openid", size, "tw-mr-2")
}

// CreateGothProvider creates a GothProvider from this Provider
func (o *OpenIDProvider) CreateGothProvider(providerName, callbackURL string, source *Source) (goth.Provider, error) {
	scopes := setting.OAuth2Client.OpenIDConnectScopes
	if len(scopes) == 0 {
		scopes = append(scopes, source.Scopes...)
	}

	provider, err := openidConnect.New(source.ClientID, source.ClientSecret, callbackURL, source.OpenIDConnectAutoDiscoveryURL, scopes...)
	if err != nil {
		log.Warn("Failed to create OpenID Connect Provider with name '%s' with url '%s': %v", providerName, source.OpenIDConnectAutoDiscoveryURL, err)
	}
	return provider, err
}

// CustomURLSettings returns the custom url settings for this provider
func (o *OpenIDProvider) CustomURLSettings() *CustomURLSettings {
	return nil
}

var _ GothProvider = &OpenIDProvider{}

func init() {
	RegisterGothProvider(&OpenIDProvider{})
}
