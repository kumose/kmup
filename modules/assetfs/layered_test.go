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

package assetfs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayered(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "assetfs-layers")
	dir1 := filepath.Join(dir, "l1")
	dir2 := filepath.Join(dir, "l2")

	mkdir := func(elems ...string) {
		assert.NoError(t, os.MkdirAll(filepath.Join(elems...), 0o755))
	}
	write := func(content string, elems ...string) {
		assert.NoError(t, os.WriteFile(filepath.Join(elems...), []byte(content), 0o644))
	}

	// d1 & f1: only in "l1"; d2 & f2: only in "l2"
	// da & fa: in both "l1" and "l2"
	mkdir(dir1, "d1")
	mkdir(dir1, "da")
	mkdir(dir1, "da/sub1")

	mkdir(dir2, "d2")
	mkdir(dir2, "da")
	mkdir(dir2, "da/sub2")

	write("dummy", dir1, ".DS_Store")
	write("f1", dir1, "f1")
	write("fa-1", dir1, "fa")
	write("d1-f", dir1, "d1/f")
	write("da-f-1", dir1, "da/f")

	write("f2", dir2, "f2")
	write("fa-2", dir2, "fa")
	write("d2-f", dir2, "d2/f")
	write("da-f-2", dir2, "da/f")

	assets := Layered(Local("l1", dir1), Local("l2", dir2))

	f, err := assets.Open("f1")
	assert.NoError(t, err)
	bs, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, "f1", string(bs))
	_ = f.Close()

	assertRead := func(expected string, expectedErr error, elems ...string) {
		bs, err := assets.ReadFile(elems...)
		if err != nil {
			assert.ErrorIs(t, err, expectedErr)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, expected, string(bs))
		}
	}
	assertRead("f1", nil, "f1")
	assertRead("f2", nil, "f2")
	assertRead("fa-1", nil, "fa")

	assertRead("d1-f", nil, "d1/f")
	assertRead("d2-f", nil, "d2/f")
	assertRead("da-f-1", nil, "da/f")

	assertRead("", fs.ErrNotExist, "no-such")

	files, err := assets.ListFiles(".", true)
	assert.NoError(t, err)
	assert.Equal(t, []string{"f1", "f2", "fa"}, files)

	files, err = assets.ListFiles(".", false)
	assert.NoError(t, err)
	assert.Equal(t, []string{"d1", "d2", "da"}, files)

	files, err = assets.ListFiles(".")
	assert.NoError(t, err)
	assert.Equal(t, []string{"d1", "d2", "da", "f1", "f2", "fa"}, files)

	files, err = assets.ListAllFiles(".", true)
	assert.NoError(t, err)
	assert.Equal(t, []string{"d1/f", "d2/f", "da/f", "f1", "f2", "fa"}, files)

	files, err = assets.ListAllFiles(".", false)
	assert.NoError(t, err)
	assert.Equal(t, []string{"d1", "d2", "da", "da/sub1", "da/sub2"}, files)

	files, err = assets.ListAllFiles(".")
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"d1", "d1/f",
		"d2", "d2/f",
		"da", "da/f", "da/sub1", "da/sub2",
		"f1", "f2", "fa",
	}, files)

	assert.Empty(t, assets.GetFileLayerName("no-such"))
	assert.Equal(t, "l1", assets.GetFileLayerName("f1"))
	assert.Equal(t, "l2", assets.GetFileLayerName("f2"))
}
