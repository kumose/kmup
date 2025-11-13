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

package markup

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestRenderHelperCodePreview(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	ctx, _ := contexttest.MockContext(t, "/", contexttest.MockContextOption{Render: templates.HTMLRenderer()})
	htm, err := renderRepoFileCodePreview(ctx, markup.RenderCodePreviewOptions{
		FullURL:   "http://full",
		OwnerName: "user2",
		RepoName:  "repo1",
		CommitID:  "65f1bf27bc3bf70f64657658635e66094edbcb4d",
		FilePath:  "README.md",
		LineStart: 1,
		LineStop:  2,
	})
	assert.NoError(t, err)
	assert.Equal(t, `<div class="code-preview-container file-content">
	<div class="code-preview-header">
		<a href="http://full" class="tw-font-semibold" rel="nofollow">repo1/README.md</a>
		repo.code_preview_line_from_to:1,2,<a href="/user2/repo1/commit/65f1bf27bc3bf70f64657658635e66094edbcb4d" class="muted tw-font-mono tw-text-text" rel="nofollow">65f1bf27bc</a>
	</div>
	<table class="file-view">
		<tbody><tr>
				<td class="lines-num"><span data-line-number="1"></span></td>
				<td class="lines-code chroma"><div class="code-inner"><span class="gh"># repo1</div></td>
			</tr><tr>
				<td class="lines-num"><span data-line-number="2"></span></td>
				<td class="lines-code chroma"><div class="code-inner"></span><span class="gh"></span></div></td>
			</tr></tbody>
	</table>
</div>
`, string(htm))

	ctx, _ = contexttest.MockContext(t, "/", contexttest.MockContextOption{Render: templates.HTMLRenderer()})
	htm, err = renderRepoFileCodePreview(ctx, markup.RenderCodePreviewOptions{
		FullURL:   "http://full",
		OwnerName: "user2",
		RepoName:  "repo1",
		CommitID:  "65f1bf27bc3bf70f64657658635e66094edbcb4d",
		FilePath:  "README.md",
		LineStart: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, `<div class="code-preview-container file-content">
	<div class="code-preview-header">
		<a href="http://full" class="tw-font-semibold" rel="nofollow">repo1/README.md</a>
		repo.code_preview_line_in:1,<a href="/user2/repo1/commit/65f1bf27bc3bf70f64657658635e66094edbcb4d" class="muted tw-font-mono tw-text-text" rel="nofollow">65f1bf27bc</a>
	</div>
	<table class="file-view">
		<tbody><tr>
				<td class="lines-num"><span data-line-number="1"></span></td>
				<td class="lines-code chroma"><div class="code-inner"><span class="gh"># repo1</div></td>
			</tr></tbody>
	</table>
</div>
`, string(htm))

	ctx, _ = contexttest.MockContext(t, "/", contexttest.MockContextOption{Render: templates.HTMLRenderer()})
	_, err = renderRepoFileCodePreview(ctx, markup.RenderCodePreviewOptions{
		FullURL:   "http://full",
		OwnerName: "user15",
		RepoName:  "big_test_private_1",
		CommitID:  "65f1bf27bc3bf70f64657658635e66094edbcb4d",
		FilePath:  "README.md",
		LineStart: 1,
		LineStop:  10,
	})
	assert.ErrorIs(t, err, util.ErrPermissionDenied)
}
