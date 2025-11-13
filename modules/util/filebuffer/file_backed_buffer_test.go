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

package filebuffer

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBackedBuffer(t *testing.T) {
	cases := []struct {
		MaxMemorySize int
		Data          string
	}{
		{5, "test"},
		{5, "testtest"},
	}

	for _, c := range cases {
		buf := New(c.MaxMemorySize, t.TempDir())
		_, err := io.Copy(buf, strings.NewReader(c.Data))
		assert.NoError(t, err)

		assert.EqualValues(t, len(c.Data), buf.Size())

		data, err := io.ReadAll(buf)
		assert.NoError(t, err)
		assert.Equal(t, c.Data, string(data))

		assert.NoError(t, buf.Close())
	}
}
