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

package hash

import (
	"encoding/hex"
)

// DummyHasher implements PasswordHasher and is a dummy hasher that simply
// puts the password in place with its salt
// This SHOULD NOT be used in production and is provided to make the integration
// tests faster only
type DummyHasher struct{}

// HashWithSaltBytes a provided password and salt
func (hasher *DummyHasher) HashWithSaltBytes(password string, salt []byte) string {
	if hasher == nil {
		return ""
	}

	if len(salt) == 10 {
		return string(salt) + ":" + password
	}

	return hex.EncodeToString(salt) + ":" + password
}

// NewDummyHasher is a factory method to create a DummyHasher
// Any provided configuration is ignored
func NewDummyHasher(_ string) *DummyHasher {
	return &DummyHasher{}
}
