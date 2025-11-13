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

func TestHashFilePathForWebUI(t *testing.T) {
	assert.Equal(t,
		"8843d7f92416211de9ebb963ff4ce28125932878",
		HashFilePathForWebUI("foobar"),
	)
}

func TestSplitCommitTitleBody(t *testing.T) {
	title, body := SplitCommitTitleBody("啊bcdefg", 4)
	assert.Equal(t, "啊…", title)
	assert.Equal(t, "…bcdefg", body)

	title, body = SplitCommitTitleBody("abcdefg\n1234567", 4)
	assert.Equal(t, "a…", title)
	assert.Equal(t, "…bcdefg\n1234567", body)

	title, body = SplitCommitTitleBody("abcdefg\n1234567", 100)
	assert.Equal(t, "abcdefg", title)
	assert.Equal(t, "1234567", body)
}
