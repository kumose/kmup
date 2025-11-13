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

package storage

import (
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func testStorageIterator(t *testing.T, typStr Type, cfg *setting.Storage) {
	l, err := NewStorage(typStr, cfg)
	assert.NoError(t, err)

	testFiles := [][]string{
		{"a/1.txt", "a1"},
		{"/a/1.txt", "aa1"}, // same as above, but with leading slash that will be trim
		{"ab/1.txt", "ab1"},
		{"b/1.txt", "b1"},
		{"b/2.txt", "b2"},
		{"b/3.txt", "b3"},
		{"b/x 4.txt", "bx4"},
	}
	for _, f := range testFiles {
		_, err = l.Save(f[0], strings.NewReader(f[1]), -1)
		assert.NoError(t, err)
	}

	expectedList := map[string][]string{
		"a":           {"a/1.txt"},
		"b":           {"b/1.txt", "b/2.txt", "b/3.txt", "b/x 4.txt"},
		"":            {"a/1.txt", "b/1.txt", "b/2.txt", "b/3.txt", "b/x 4.txt", "ab/1.txt"},
		"/":           {"a/1.txt", "b/1.txt", "b/2.txt", "b/3.txt", "b/x 4.txt", "ab/1.txt"},
		".":           {"a/1.txt", "b/1.txt", "b/2.txt", "b/3.txt", "b/x 4.txt", "ab/1.txt"},
		"a/b/../../a": {"a/1.txt"},
	}
	for dir, expected := range expectedList {
		count := 0
		err = l.IterateObjects(dir, func(path string, f Object) error {
			defer f.Close()
			assert.Contains(t, expected, path)
			count++
			return nil
		})
		assert.NoError(t, err)
		assert.Len(t, expected, count)
	}
}
