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

package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadMarkup(t *testing.T) {
	cfg, _ := NewConfigProviderFromData(``)
	loadMarkupFrom(cfg)
	assert.Equal(t, MarkdownMathCodeBlockOptions{ParseInlineDollar: true, ParseBlockDollar: true}, Markdown.MathCodeBlockOptions)
	assert.Equal(t, MarkdownRenderOptions{NewLineHardBreak: true, ShortIssuePattern: true}, Markdown.RenderOptionsComment)
	assert.Equal(t, MarkdownRenderOptions{ShortIssuePattern: true}, Markdown.RenderOptionsWiki)
	assert.Equal(t, MarkdownRenderOptions{}, Markdown.RenderOptionsRepoFile)

	t.Run("Math", func(t *testing.T) {
		cfg, _ = NewConfigProviderFromData(`
[markdown]
MATH_CODE_BLOCK_DETECTION = none
`)
		loadMarkupFrom(cfg)
		assert.Equal(t, MarkdownMathCodeBlockOptions{}, Markdown.MathCodeBlockOptions)

		cfg, _ = NewConfigProviderFromData(`
[markdown]
MATH_CODE_BLOCK_DETECTION = inline-dollar, inline-parentheses, block-dollar, block-square-brackets
`)
		loadMarkupFrom(cfg)
		assert.Equal(t, MarkdownMathCodeBlockOptions{ParseInlineDollar: true, ParseInlineParentheses: true, ParseBlockDollar: true, ParseBlockSquareBrackets: true}, Markdown.MathCodeBlockOptions)
	})

	t.Run("Render", func(t *testing.T) {
		cfg, _ = NewConfigProviderFromData(`
[markdown]
RENDER_OPTIONS_COMMENT = none
`)
		loadMarkupFrom(cfg)
		assert.Equal(t, MarkdownRenderOptions{}, Markdown.RenderOptionsComment)

		cfg, _ = NewConfigProviderFromData(`
[markdown]
RENDER_OPTIONS_REPO_FILE = short-issue-pattern, new-line-hard-break
`)
		loadMarkupFrom(cfg)
		assert.Equal(t, MarkdownRenderOptions{NewLineHardBreak: true, ShortIssuePattern: true}, Markdown.RenderOptionsRepoFile)
	})
}
