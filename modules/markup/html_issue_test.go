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

	"github.com/kumose/kmup/modules/htmlutil"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/markup/markdown"
	testModule "github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_IssueList(t *testing.T) {
	defer testModule.MockVariableValue(&markup.RenderBehaviorForTesting.DisableAdditionalAttributes, true)()
	markup.Init(&markup.RenderHelperFuncs{
		RenderRepoIssueIconTitle: func(ctx context.Context, opts markup.RenderIssueIconTitleOptions) (template.HTML, error) {
			return htmlutil.HTMLFormat("<div>issue #%d</div>", opts.IssueIndex), nil
		},
	})

	test := func(input, expected string) {
		rctx := markup.NewTestRenderContext(markup.TestAppURL, map[string]string{
			"user": "test-user", "repo": "test-repo",
			"markupAllowShortIssuePattern": "true",
			"footnoteContextId":            "12345",
		})
		out, err := markdown.RenderString(rctx, input)
		require.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
	}

	t.Run("NormalIssueRef", func(t *testing.T) {
		test(
			"#12345",
			`<p><a href="/test-user/test-repo/issues/12345" class="ref-issue" rel="nofollow">#12345</a></p>`,
		)
	})

	t.Run("ListIssueRef", func(t *testing.T) {
		test(
			"* #12345",
			`<ul>
<li><div>issue #12345</div></li>
</ul>`,
		)
	})

	t.Run("ListIssueRefNormal", func(t *testing.T) {
		test(
			"* foo #12345 bar",
			`<ul>
<li>foo <a href="/test-user/test-repo/issues/12345" class="ref-issue" rel="nofollow">#12345</a> bar</li>
</ul>`,
		)
	})

	t.Run("ListTodoIssueRef", func(t *testing.T) {
		test(
			"* [ ] #12345",
			`<ul>
<li class="task-list-item"><input type="checkbox" disabled="" data-source-position="2"/><div>issue #12345</div></li>
</ul>`,
		)
	})

	t.Run("IssueFootnote", func(t *testing.T) {
		test(
			"foo[^1][^2]\n\n[^1]: bar\n[^2]: baz",
			`<p>foo<sup id="fnref:user-content-1-12345"><a href="#fn:user-content-1-12345" rel="nofollow">1 </a></sup><sup id="fnref:user-content-2-12345"><a href="#fn:user-content-2-12345" rel="nofollow">2 </a></sup></p>
<div>
<hr/>
<ol>
<li id="fn:user-content-1-12345">
<p>bar <a href="#fnref:user-content-1-12345" rel="nofollow">↩︎</a></p>
</li>
<li id="fn:user-content-2-12345">
<p>baz <a href="#fnref:user-content-2-12345" rel="nofollow">↩︎</a></p>
</li>
</ol>
</div>`,
		)
	})
}
