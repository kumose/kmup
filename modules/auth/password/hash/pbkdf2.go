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
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/kumose/kmup/modules/log"

	"golang.org/x/crypto/pbkdf2"
)

func init() {
	MustRegister("pbkdf2", NewPBKDF2Hasher)
}

// PBKDF2Hasher implements PasswordHasher
// and uses the PBKDF2 key derivation function.
type PBKDF2Hasher struct {
	iter, keyLen int
}

// HashWithSaltBytes a provided password and salt
func (hasher *PBKDF2Hasher) HashWithSaltBytes(password string, salt []byte) string {
	if hasher == nil {
		return ""
	}
	return hex.EncodeToString(pbkdf2.Key([]byte(password), salt, hasher.iter, hasher.keyLen, sha256.New))
}

// NewPBKDF2Hasher is a factory method to create an PBKDF2Hasher
// config should be either empty or of the form:
// "<iter>$<keyLen>", where <x> is the string representation
// of an integer
func NewPBKDF2Hasher(config string) *PBKDF2Hasher {
	// This default configuration uses the following parameters:
	// iter=10000, keyLen=50.
	// This matches the original configuration for `pbkdf2` prior to storing parameters
	// in the database.
	// THESE VALUES MUST NOT BE CHANGED OR BACKWARDS COMPATIBILITY WILL BREAK
	hasher := &PBKDF2Hasher{
		iter:   10_000,
		keyLen: 50,
	}

	if config == "" {
		return hasher
	}

	vals := strings.SplitN(config, "$", 2)
	if len(vals) != 2 {
		log.Error("invalid pbkdf2 hash spec %s", config)
		return nil
	}

	var err error
	hasher.iter, err = parseIntParam(vals[0], "iter", "pbkdf2", config, nil)
	hasher.keyLen, err = parseIntParam(vals[1], "keyLen", "pbkdf2", config, err)
	if err != nil {
		return nil
	}

	return hasher
}
