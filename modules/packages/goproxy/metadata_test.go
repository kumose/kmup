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

package goproxy

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	packageName    = "kmup.com/kumose/kmup"
	packageVersion = "v0.0.1"
)

func TestParsePackage(t *testing.T) {
	createArchive := func(files map[string][]byte) *bytes.Reader {
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		for name, content := range files {
			w, _ := zw.Create(name)
			w.Write(content)
		}
		zw.Close()
		return bytes.NewReader(buf.Bytes())
	}

	t.Run("EmptyPackage", func(t *testing.T) {
		data := createArchive(nil)

		p, err := ParsePackage(data, int64(data.Len()))
		assert.Nil(t, p)
		assert.ErrorIs(t, err, ErrInvalidStructure)
	})

	t.Run("InvalidNameOrVersionStructure", func(t *testing.T) {
		data := createArchive(map[string][]byte{
			packageName + "/" + packageVersion + "/go.mod": {},
		})

		p, err := ParsePackage(data, int64(data.Len()))
		assert.Nil(t, p)
		assert.ErrorIs(t, err, ErrInvalidStructure)
	})

	t.Run("GoModFileInWrongDirectory", func(t *testing.T) {
		data := createArchive(map[string][]byte{
			packageName + "@" + packageVersion + "/subdir/go.mod": {},
		})

		p, err := ParsePackage(data, int64(data.Len()))
		assert.NotNil(t, p)
		assert.NoError(t, err)
		assert.Equal(t, packageName, p.Name)
		assert.Equal(t, packageVersion, p.Version)
		assert.Equal(t, "module kmup.com/kumose/kmup", p.GoMod)
	})

	t.Run("Valid", func(t *testing.T) {
		data := createArchive(map[string][]byte{
			packageName + "@" + packageVersion + "/subdir/go.mod": []byte("invalid"),
			packageName + "@" + packageVersion + "/go.mod":        []byte("valid"),
		})

		p, err := ParsePackage(data, int64(data.Len()))
		assert.NotNil(t, p)
		assert.NoError(t, err)
		assert.Equal(t, packageName, p.Name)
		assert.Equal(t, packageVersion, p.Version)
		assert.Equal(t, "valid", p.GoMod)
	})
}
