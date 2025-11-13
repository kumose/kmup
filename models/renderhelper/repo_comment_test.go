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

func TestRepoComment(t *testing.T) {
	unittest.PrepareTestEnv(t)

	repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	t.Run("AutoLink", func(t *testing.T) {
		rctx := NewRenderContextRepoComment(t.Context(), repo1).WithMarkupType(markdown.MarkupName)
		rendered, err := markup.RenderString(rctx, `
65f1bf27bc3bf70f64657658635e66094edbcb4d
#1
@user2
`)
		assert.NoError(t, err)
		assert.Equal(t,
			`<p><a href="/user2/repo1/commit/65f1bf27bc3bf70f64657658635e66094edbcb4d" rel="nofollow"><code>65f1bf27bc</code></a><br/>
<a href="/user2/repo1/issues/1" class="ref-issue" rel="nofollow">#1</a><br/>
<a href="/user2" rel="nofollow">@user2</a></p>
`, rendered)
	})

	t.Run("AbsoluteAndRelative", func(t *testing.T) {
		rctx := NewRenderContextRepoComment(t.Context(), repo1).WithMarkupType(markdown.MarkupName)

		// It is Kmup's old behavior, the relative path is resolved to the repo path
		// It is different from GitHub, GitHub resolves relative links to current page's path
		rendered, err := markup.RenderString(rctx, `
[/test](/test)
[./test](./test)
![/image](/image)
![./image](./image)
`)
		assert.NoError(t, err)
		assert.Equal(t,
			`<p><a href="/user2/repo1/test" rel="nofollow">/test</a><br/>
<a href="/user2/repo1/test" rel="nofollow">./test</a><br/>
<a href="/user2/repo1/image" target="_blank" rel="nofollow noopener"><img src="/user2/repo1/image" alt="/image"/></a><br/>
<a href="/user2/repo1/image" target="_blank" rel="nofollow noopener"><img src="/user2/repo1/image" alt="./image"/></a></p>
`, rendered)
	})

	t.Run("WithCurrentRefPath", func(t *testing.T) {
		rctx := NewRenderContextRepoComment(t.Context(), repo1, RepoCommentOptions{CurrentRefPath: "/commit/1234"}).
			WithMarkupType(markdown.MarkupName)

		// the ref path is only used to render commit message: a commit message is rendered at the commit page with its commit ID path
		rendered, err := markup.RenderString(rctx, `
[/test](/test)
[./test](./test)
![/image](/image)
![./image](./image)
`)
		assert.NoError(t, err)
		assert.Equal(t, `<p><a href="/user2/repo1/test" rel="nofollow">/test</a><br/>
<a href="/user2/repo1/commit/1234/test" rel="nofollow">./test</a><br/>
<a href="/user2/repo1/image" target="_blank" rel="nofollow noopener"><img src="/user2/repo1/image" alt="/image"/></a><br/>
<a href="/user2/repo1/commit/1234/image" target="_blank" rel="nofollow noopener"><img src="/user2/repo1/commit/1234/image" alt="./image"/></a></p>
`, rendered)
	})

	t.Run("NoRepo", func(t *testing.T) {
		rctx := NewRenderContextRepoComment(t.Context(), nil).WithMarkupType(markdown.MarkupName)
		rendered, err := markup.RenderString(rctx, "any")
		assert.NoError(t, err)
		assert.Equal(t, "<p>any</p>\n", rendered)
	})
}
