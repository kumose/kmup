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

package smtp

import (
	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/json"
)

//   _________   __________________________
//  /   _____/  /     \__    ___/\______   \
//  \_____  \  /  \ /  \|    |    |     ___/
//  /        \/    Y    \    |    |    |
// /_______  /\____|__  /____|    |____|
//         \/         \/

// Source holds configuration for the SMTP login source.
type Source struct {
	auth.ConfigBase `json:"-"`

	Auth           string
	Host           string
	Port           int
	AllowedDomains string `xorm:"TEXT"`
	ForceSMTPS     bool
	SkipVerify     bool
	HeloHostname   string
	DisableHelo    bool
}

// FromDB fills up an SMTPConfig from serialized format.
func (source *Source) FromDB(bs []byte) error {
	return json.UnmarshalHandleDoubleEncode(bs, &source)
}

// ToDB exports an SMTPConfig to a serialized format.
func (source *Source) ToDB() ([]byte, error) {
	return json.Marshal(source)
}

// IsSkipVerify returns if SkipVerify is set
func (source *Source) IsSkipVerify() bool {
	return source.SkipVerify
}

// HasTLS returns true for SMTP
func (source *Source) HasTLS() bool {
	return true
}

// UseTLS returns if TLS is set
func (source *Source) UseTLS() bool {
	return source.ForceSMTPS || source.Port == 465
}

func init() {
	auth.RegisterTypeConfig(auth.SMTP, &Source{})
}
