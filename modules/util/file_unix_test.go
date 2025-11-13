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

//go:build !windows

package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyUmask(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "test-filemode-")
	assert.NoError(t, err)

	err = os.Chmod(f.Name(), 0o777)
	assert.NoError(t, err)
	st, err := os.Stat(f.Name())
	assert.NoError(t, err)
	assert.EqualValues(t, 0o777, st.Mode().Perm()&0o777)

	oldDefaultUmask := defaultUmask
	defaultUmask = 0o037
	defer func() {
		defaultUmask = oldDefaultUmask
	}()
	err = ApplyUmask(f.Name(), os.ModePerm)
	assert.NoError(t, err)
	st, err = os.Stat(f.Name())
	assert.NoError(t, err)
	assert.EqualValues(t, 0o740, st.Mode().Perm()&0o777)
}
