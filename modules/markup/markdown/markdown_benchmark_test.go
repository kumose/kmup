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

package markdown_test

import (
	"testing"

	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/markup/markdown"
)

func BenchmarkSpecializedMarkdown(b *testing.B) {
	// 240856	      4719 ns/op
	for b.Loop() {
		markdown.SpecializedMarkdown(&markup.RenderContext{})
	}
}

func BenchmarkMarkdownRender(b *testing.B) {
	// 23202	     50840 ns/op
	for b.Loop() {
		_, _ = markdown.RenderString(markup.NewTestRenderContext(), "https://example.com\n- a\n- b\n")
	}
}
