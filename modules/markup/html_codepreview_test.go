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

package markup_test

import (
	"context"
	"html/template"
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/markup/markdown"

	"github.com/stretchr/testify/assert"
)

func TestRenderCodePreview(t *testing.T) {
	markup.Init(&markup.RenderHelperFuncs{
		RenderRepoFileCodePreview: func(ctx context.Context, options markup.RenderCodePreviewOptions) (template.HTML, error) {
			return "<div>code preview</div>", nil
		},
	})
	test := func(input, expected string) {
		buffer, err := markup.RenderString(markup.NewTestRenderContext().WithMarkupType(markdown.MarkupName), input)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(buffer))
	}
	test("http://localhost:3326/owner/repo/src/commit/0123456789/foo/bar.md#L10-L20", "<p><div>code preview</div></p>")
	test("http://other/owner/repo/src/commit/0123456789/foo/bar.md#L10-L20", `<p><a href="http://other/owner/repo/src/commit/0123456789/foo/bar.md#L10-L20" rel="nofollow">http://other/owner/repo/src/commit/0123456789/foo/bar.md#L10-L20</a></p>`)
}
