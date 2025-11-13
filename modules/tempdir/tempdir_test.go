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

package tempdir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTempDir(t *testing.T) {
	base := t.TempDir()

	t.Run("Create", func(t *testing.T) {
		td := New(base, "sub1/sub2") // make sure the sub dir supports "/" in the path
		assert.Equal(t, filepath.Join(base, "sub1", "sub2"), td.JoinPath())
		assert.Equal(t, filepath.Join(base, "sub1", "sub2/test"), td.JoinPath("test"))

		t.Run("MkdirTempRandom", func(t *testing.T) {
			s, cleanup, err := td.MkdirTempRandom("foo")
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(s, filepath.Join(base, "sub1/sub2", "foo")))

			_, err = os.Stat(s)
			assert.NoError(t, err)
			cleanup()
			_, err = os.Stat(s)
			assert.ErrorIs(t, err, os.ErrNotExist)
		})

		t.Run("CreateTempFileRandom", func(t *testing.T) {
			f, cleanup, err := td.CreateTempFileRandom("foo", "bar")
			filename := f.Name()
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(filename, filepath.Join(base, "sub1/sub2", "foo", "bar")))
			_, err = os.Stat(filename)
			assert.NoError(t, err)
			cleanup()
			_, err = os.Stat(filename)
			assert.ErrorIs(t, err, os.ErrNotExist)
		})

		t.Run("RemoveOutDated", func(t *testing.T) {
			fa1, _, err := td.CreateTempFileRandom("dir-a", "f1")
			assert.NoError(t, err)
			fa2, _, err := td.CreateTempFileRandom("dir-a", "f2")
			assert.NoError(t, err)
			_ = os.Chtimes(fa2.Name(), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour))
			fb1, _, err := td.CreateTempFileRandom("dir-b", "f1")
			assert.NoError(t, err)
			_ = os.Chtimes(fb1.Name(), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour))
			_, _, _ = fa1.Close(), fa2.Close(), fb1.Close()

			td.RemoveOutdated(time.Minute)

			_, err = os.Stat(fa1.Name())
			assert.NoError(t, err)
			_, err = os.Stat(fa2.Name())
			assert.ErrorIs(t, err, os.ErrNotExist)
			_, err = os.Stat(fb1.Name())
			assert.ErrorIs(t, err, os.ErrNotExist)
		})
	})

	t.Run("BaseNotExist", func(t *testing.T) {
		td := New(filepath.Join(base, "not-exist"), "sub")
		_, _, err := td.MkdirTempRandom("foo")
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}
