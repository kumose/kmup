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
	"golang.org/x/crypto/bcrypt"
)

func init() {
	MustRegister("bcrypt", NewBcryptHasher)
}

// BcryptHasher implements PasswordHasher
// and uses the bcrypt password hash function.
type BcryptHasher struct {
	cost int
}

// HashWithSaltBytes a provided password and salt
func (hasher *BcryptHasher) HashWithSaltBytes(password string, salt []byte) string {
	if hasher == nil {
		return ""
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), hasher.cost)
	return string(hashedPassword)
}

func (hasher *BcryptHasher) VerifyPassword(password, hashedPassword, salt string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// NewBcryptHasher is a factory method to create an BcryptHasher
// The provided config should be either empty or the string representation of the "<cost>"
// as an integer
func NewBcryptHasher(config string) *BcryptHasher {
	// This matches the original configuration for `bcrypt` prior to storing hash parameters
	// in the database.
	// THESE VALUES MUST NOT BE CHANGED OR BACKWARDS COMPATIBILITY WILL BREAK
	hasher := &BcryptHasher{
		cost: 10, // cost=10. i.e. 2^10 rounds of key expansion.
	}

	if config == "" {
		return hasher
	}
	var err error
	hasher.cost, err = parseIntParam(config, "cost", "bcrypt", config, nil)
	if err != nil {
		return nil
	}

	return hasher
}
