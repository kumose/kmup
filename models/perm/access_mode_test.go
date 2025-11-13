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

package perm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessMode(t *testing.T) {
	names := []string{"none", "read", "write", "admin"}
	for i, name := range names {
		m := ParseAccessMode(name)
		assert.Equal(t, AccessMode(i), m)
	}
	assert.Equal(t, AccessModeOwner, AccessMode(4))
	assert.Equal(t, "owner", AccessModeOwner.ToString())
	assert.Equal(t, AccessModeNone, ParseAccessMode("owner"))
	assert.Equal(t, AccessModeNone, ParseAccessMode("invalid"))
}
