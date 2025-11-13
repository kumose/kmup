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

package ldap

import (
	"strings"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/secret"
	"github.com/kumose/kmup/modules/setting"
)

// .____     ________      _____ __________
// |    |    \______ \    /  _  \\______   \
// |    |     |    |  \  /  /_\  \|     ___/
// |    |___  |    `   \/    |    \    |
// |_______ \/_______  /\____|__  /____|
//         \/        \/         \/

// Package ldap provide functions & structure to query a LDAP ldap directory
// For now, it's mainly tested again an MS Active Directory service, see README.md for more information

// Source Basic LDAP authentication service
type Source struct {
	auth.ConfigBase `json:"-"`

	Name                  string // canonical name (ie. corporate.ad)
	Host                  string // LDAP host
	Port                  int    // port number
	SecurityProtocol      SecurityProtocol
	SkipVerify            bool
	BindDN                string // DN to bind with
	BindPasswordEncrypt   string // Encrypted Bind BN password
	BindPassword          string // Bind DN password
	UserBase              string // Base search path for users
	UserDN                string // Template for the DN of the user for simple auth
	AttributeUsername     string // Username attribute
	AttributeName         string // First name attribute
	AttributeSurname      string // Surname attribute
	AttributeMail         string // E-mail attribute
	AttributesInBind      bool   // fetch attributes in bind context (not user)
	AttributeSSHPublicKey string // LDAP SSH Public Key attribute
	AttributeAvatar       string
	SearchPageSize        uint32 // Search with paging page size
	Filter                string // Query filter to validate entry
	AdminFilter           string // Query filter to check if user is admin
	RestrictedFilter      string // Query filter to check if user is restricted
	Enabled               bool   // if this source is disabled
	AllowDeactivateAll    bool   // Allow an empty search response to deactivate all users from this source
	GroupsEnabled         bool   // if the group checking is enabled
	GroupDN               string // Group Search Base
	GroupFilter           string // Group Name Filter
	GroupMemberUID        string // Group Attribute containing array of UserUID
	GroupTeamMap          string // Map LDAP groups to teams
	GroupTeamMapRemoval   bool   // Remove user from teams which are synchronized and user is not a member of the corresponding LDAP group
	UserUID               string // User Attribute listed in Group
}

// FromDB fills up a LDAPConfig from serialized format.
func (source *Source) FromDB(bs []byte) error {
	err := json.UnmarshalHandleDoubleEncode(bs, &source)
	if err != nil {
		return err
	}
	if source.BindPasswordEncrypt != "" {
		source.BindPassword, err = secret.DecryptSecret(setting.SecretKey, source.BindPasswordEncrypt)
		if err != nil {
			log.Error("Unable to decrypt bind password for LDAP source, maybe SECRET_KEY is wrong: %v", err)
		}
		source.BindPasswordEncrypt = ""
	}
	return nil
}

// ToDB exports a LDAPConfig to a serialized format.
func (source *Source) ToDB() ([]byte, error) {
	var err error
	source.BindPasswordEncrypt, err = secret.EncryptSecret(setting.SecretKey, source.BindPassword)
	if err != nil {
		return nil, err
	}
	source.BindPassword = ""
	return json.Marshal(source)
}

// SecurityProtocolName returns the name of configured security
// protocol.
func (source *Source) SecurityProtocolName() string {
	return SecurityProtocolNames[source.SecurityProtocol]
}

// IsSkipVerify returns if SkipVerify is set
func (source *Source) IsSkipVerify() bool {
	return source.SkipVerify
}

// HasTLS returns if HasTLS
func (source *Source) HasTLS() bool {
	return source.SecurityProtocol > SecurityProtocolUnencrypted
}

// UseTLS returns if UseTLS
func (source *Source) UseTLS() bool {
	return source.SecurityProtocol != SecurityProtocolUnencrypted
}

// ProvidesSSHKeys returns if this source provides SSH Keys
func (source *Source) ProvidesSSHKeys() bool {
	return strings.TrimSpace(source.AttributeSSHPublicKey) != ""
}

func init() {
	auth.RegisterTypeConfig(auth.LDAP, &Source{})
	auth.RegisterTypeConfig(auth.DLDAP, &Source{})
}
