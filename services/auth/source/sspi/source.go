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

package sspi

import (
	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/json"
)

//   _________ ___________________.___
//  /   _____//   _____/\______   \   |
//  \_____  \ \_____  \  |     ___/   |
//  /        \/        \ |    |   |   |
// /_______  /_______  / |____|   |___|
//         \/        \/

// Source holds configuration for SSPI single sign-on.
type Source struct {
	auth.ConfigBase `json:"-"`

	AutoCreateUsers      bool
	AutoActivateUsers    bool
	StripDomainNames     bool
	SeparatorReplacement string
	DefaultLanguage      string
}

// FromDB fills up an SSPIConfig from serialized format.
func (cfg *Source) FromDB(bs []byte) error {
	return json.UnmarshalHandleDoubleEncode(bs, &cfg)
}

// ToDB exports an SSPIConfig to a serialized format.
func (cfg *Source) ToDB() ([]byte, error) {
	return json.Marshal(cfg)
}

func init() {
	auth.RegisterTypeConfig(auth.SSPI, &Source{})
}
