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

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidSHAPattern(t *testing.T) {
	h := Sha1ObjectFormat
	assert.True(t, h.IsValid("fee1"))
	assert.True(t, h.IsValid("abc000"))
	assert.True(t, h.IsValid("9023902390239023902390239023902390239023"))
	assert.False(t, h.IsValid("90239023902390239023902390239023902390239023"))
	assert.False(t, h.IsValid("abc"))
	assert.False(t, h.IsValid("123g"))
	assert.False(t, h.IsValid("some random text"))
	assert.Equal(t, "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391", ComputeBlobHash(Sha1ObjectFormat, nil).String())
	assert.Equal(t, "2e65efe2a145dda7ee51d1741299f848e5bf752e", ComputeBlobHash(Sha1ObjectFormat, []byte("a")).String())
	assert.Equal(t, "473a0f4c3be8a93681a267e3b1e9a7dcda1185436fe141f7749120a303721813", ComputeBlobHash(Sha256ObjectFormat, nil).String())
	assert.Equal(t, "eb337bcee2061c5313c9a1392116b6c76039e9e30d71467ae359b36277e17dc7", ComputeBlobHash(Sha256ObjectFormat, []byte("a")).String())
}
