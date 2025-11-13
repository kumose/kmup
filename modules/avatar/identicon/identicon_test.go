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

//go:build test_avatar_identicon

package identicon

import (
	"image/color"
	"image/png"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	dir, _ := os.Getwd()
	dir = dir + "/testdata"
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		t.Errorf("can not save generated images to %s", dir)
	}

	backColor := color.White
	imgMaker, err := New(64, backColor, DarkColors...)
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		s := strconv.Itoa(i)
		img := imgMaker.Make([]byte(s))

		f, err := os.Create(dir + "/" + s + ".png")
		if !assert.NoError(t, err) {
			continue
		}
		defer f.Close()
		err = png.Encode(f, img)
		assert.NoError(t, err)
	}
}
