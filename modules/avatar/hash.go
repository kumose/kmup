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

package avatar

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// HashAvatar will generate a unique string, which ensures that when there's a
// different unique ID while the data is the same, it will generate a different
// output. It will generate the output according to:
// HEX(HASH(uniqueID || - || data))
// The hash being used is SHA256.
// The sole purpose of the unique ID is to generate a distinct hash Such that
// two unique IDs with the same data will have a different hash output.
// The "-" byte is important to ensure that data cannot be modified such that
// the first byte is a number, which could lead to a "collision" with the hash
// of another unique ID.
func HashAvatar(uniqueID int64, data []byte) string {
	h := sha256.New()
	h.Write([]byte(strconv.FormatInt(uniqueID, 10)))
	h.Write([]byte{'-'})
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
