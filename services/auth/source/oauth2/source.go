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
	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/json"
)

// Source holds configuration for the OAuth2 login source.
type Source struct {
	auth.ConfigBase `json:"-"`

	Provider                      string
	ClientID                      string
	ClientSecret                  string
	OpenIDConnectAutoDiscoveryURL string
	CustomURLMapping              *CustomURLMapping
	IconURL                       string

	Scopes              []string
	RequiredClaimName   string
	RequiredClaimValue  string
	GroupClaimName      string
	AdminGroup          string
	GroupTeamMap        string
	GroupTeamMapRemoval bool
	RestrictedGroup     string

	SSHPublicKeyClaimName string
	FullNameClaimName     string
}

// FromDB fills up an OAuth2Config from serialized format.
func (source *Source) FromDB(bs []byte) error {
	return json.UnmarshalHandleDoubleEncode(bs, &source)
}

// ToDB exports an OAuth2Config to a serialized format.
func (source *Source) ToDB() ([]byte, error) {
	return json.Marshal(source)
}

func init() {
	auth.RegisterTypeConfig(auth.OAuth2, &Source{})
}
