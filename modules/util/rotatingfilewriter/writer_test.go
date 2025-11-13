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

package rotatingfilewriter

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressOldFile(t *testing.T) {
	tmpDir := t.TempDir()
	fname := filepath.Join(tmpDir, "test")
	nonGzip := filepath.Join(tmpDir, "test-nonGzip")

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0o660)
	assert.NoError(t, err)
	ng, err := os.OpenFile(nonGzip, os.O_CREATE|os.O_WRONLY, 0o660)
	assert.NoError(t, err)

	for range 999 {
		f.WriteString("This is a test file\n")
		ng.WriteString("This is a test file\n")
	}
	f.Close()
	ng.Close()

	err = compressOldFile(fname, gzip.DefaultCompression)
	assert.NoError(t, err)

	_, err = os.Lstat(fname + ".gz")
	assert.NoError(t, err)

	f, err = os.Open(fname + ".gz")
	assert.NoError(t, err)
	zr, err := gzip.NewReader(f)
	assert.NoError(t, err)
	data, err := io.ReadAll(zr)
	assert.NoError(t, err)
	original, err := os.ReadFile(nonGzip)
	assert.NoError(t, err)
	assert.Equal(t, original, data)
}
