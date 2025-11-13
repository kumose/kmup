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

package avatar_test

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/kumose/kmup/modules/avatar"

	"github.com/stretchr/testify/assert"
)

func Test_HashAvatar(t *testing.T) {
	myImage := image.NewRGBA(image.Rect(0, 0, 32, 32))
	var buff bytes.Buffer
	png.Encode(&buff, myImage)

	assert.Equal(t, "9ddb5bac41d57e72aa876321d0c09d71090c05f94bc625303801be2f3240d2cb", avatar.HashAvatar(1, buff.Bytes()))
	assert.Equal(t, "9a5d44e5d637b9582a976676e8f3de1dccd877c2fe3e66ca3fab1629f2f47609", avatar.HashAvatar(8, buff.Bytes()))
	assert.Equal(t, "ed7399158672088770de6f5211ce15528ebd675e92fc4fc060c025f4b2794ccb", avatar.HashAvatar(1024, buff.Bytes()))
	assert.Equal(t, "161178642c7d59eb25a61dddced5e6b66eae1c70880d5f148b1b497b767e72d9", avatar.HashAvatar(1024, []byte{}))
}
