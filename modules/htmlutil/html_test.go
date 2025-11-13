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

package htmlutil

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStringer struct{}

func (t testStringer) String() string {
	return "&StringMethod"
}

func TestHTMLFormat(t *testing.T) {
	assert.Equal(t, template.HTML("<a>&lt; < 1</a>"), HTMLFormat("<a>%s %s %d</a>", "<", template.HTML("<"), 1))
	assert.Equal(t, template.HTML("%!s(<nil>)"), HTMLFormat("%s", nil))
	assert.Equal(t, template.HTML("&lt;&gt;"), HTMLFormat("%s", template.URL("<>")))
	assert.Equal(t, template.HTML("&amp;StringMethod &amp;StringMethod"), HTMLFormat("%s %s", testStringer{}, &testStringer{}))
}
