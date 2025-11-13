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

package svg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	res := Normalize([]byte("foo"), 1)
	assert.Equal(t, "foo", string(res))

	res = Normalize([]byte(`<?xml version="1.0"?>
<!--
comment
-->
<svg xmlns = "...">content</svg>`), 1)
	assert.Equal(t, `<svg width="1" height="1" class="svg">content</svg>`, string(res))

	res = Normalize([]byte(`<svg
width="100"
class="svg-icon"
>content</svg>`), 16)

	assert.Equal(t, `<svg class="svg-icon" width="16" height="16">content</svg>`, string(res))
}
