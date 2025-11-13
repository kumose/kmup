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

package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	hex, err := EncryptSecret("foo", "baz")
	assert.NoError(t, err)
	str, _ := DecryptSecret("foo", hex)
	assert.Equal(t, "baz", str)

	hex, err = EncryptSecret("bar", "baz")
	assert.NoError(t, err)
	str, _ = DecryptSecret("foo", hex)
	assert.NotEqual(t, "baz", str)

	_, err = DecryptSecret("a", "b")
	assert.ErrorContains(t, err, "invalid hex string")

	_, err = DecryptSecret("a", "bb")
	assert.ErrorContains(t, err, "the key (maybe SECRET_KEY?) might be incorrect: AesDecrypt ciphertext too short")

	_, err = DecryptSecret("a", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	assert.ErrorContains(t, err, "the key (maybe SECRET_KEY?) might be incorrect: AesDecrypt invalid decrypted base64 string")
}
