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

package renderhelper

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/markup/markdown"

	"github.com/stretchr/testify/assert"
)

func TestSimpleDocument(t *testing.T) {
	unittest.PrepareTestEnv(t)
	rctx := NewRenderContextSimpleDocument(t.Context(), "/base").WithMarkupType(markdown.MarkupName)
	rendered, err := markup.RenderString(rctx, `
65f1bf27bc3bf70f64657658635e66094edbcb4d
#1
@user2

[/test](/test)
[./test](./test)
![/image](/image)
![./image](./image)
`)
	assert.NoError(t, err)
	assert.Equal(t,
		`<p>65f1bf27bc3bf70f64657658635e66094edbcb4d
#1
<a href="/user2" rel="nofollow">@user2</a></p>
<p><a href="/base/test" rel="nofollow">/test</a>
<a href="/base/test" rel="nofollow">./test</a>
<a href="/base/image" target="_blank" rel="nofollow noopener"><img src="/base/image" alt="/image"/></a>
<a href="/base/image" target="_blank" rel="nofollow noopener"><img src="/base/image" alt="./image"/></a></p>
`, rendered)
}
