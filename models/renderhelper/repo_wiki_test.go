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

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/markup/markdown"

	"github.com/stretchr/testify/assert"
)

func TestRepoWiki(t *testing.T) {
	unittest.PrepareTestEnv(t)
	repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	t.Run("AutoLink", func(t *testing.T) {
		rctx := NewRenderContextRepoWiki(t.Context(), repo1).WithMarkupType(markdown.MarkupName)
		rendered, err := markup.RenderString(rctx, `
65f1bf27bc3bf70f64657658635e66094edbcb4d
#1
@user2
`)
		assert.NoError(t, err)
		assert.Equal(t,
			`<p><a href="/user2/repo1/commit/65f1bf27bc3bf70f64657658635e66094edbcb4d" rel="nofollow"><code>65f1bf27bc</code></a>
<a href="/user2/repo1/issues/1" class="ref-issue" rel="nofollow">#1</a>
<a href="/user2" rel="nofollow">@user2</a></p>
`, rendered)
	})

	t.Run("AbsoluteAndRelative", func(t *testing.T) {
		rctx := NewRenderContextRepoWiki(t.Context(), repo1).WithMarkupType(markdown.MarkupName)
		rendered, err := markup.RenderString(rctx, `
[/test](/test)
[./test](./test)
![/image](/image)
![./image](./image)
`)
		assert.NoError(t, err)
		assert.Equal(t,
			`<p><a href="/user2/repo1/wiki/test" rel="nofollow">/test</a>
<a href="/user2/repo1/wiki/test" rel="nofollow">./test</a>
<a href="/user2/repo1/wiki/image" target="_blank" rel="nofollow noopener"><img src="/user2/repo1/wiki/raw/image" alt="/image"/></a>
<a href="/user2/repo1/wiki/image" target="_blank" rel="nofollow noopener"><img src="/user2/repo1/wiki/raw/image" alt="./image"/></a></p>
`, rendered)
	})

	t.Run("PathInTag", func(t *testing.T) {
		rctx := NewRenderContextRepoWiki(t.Context(), repo1).WithMarkupType(markdown.MarkupName)
		rendered, err := markup.RenderString(rctx, `
<img src="LINK">
<video src="LINK">
`)
		assert.NoError(t, err)
		assert.Equal(t, `<a href="/user2/repo1/wiki/LINK" target="_blank" rel="nofollow noopener"><img src="/user2/repo1/wiki/raw/LINK"/></a>
<video src="/user2/repo1/wiki/raw/LINK">
</video>`, rendered)
	})
}
