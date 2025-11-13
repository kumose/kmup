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

package math

import (
	"bytes"

	"github.com/kumose/kmup/modules/markup/internal"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Inline render output:
// <code class="language-math">...</code>

// InlineRenderer is an inline renderer
type InlineRenderer struct {
	renderInternal *internal.RenderInternal
}

// NewInlineRenderer returns a new renderer for inline math
func NewInlineRenderer(renderInternal *internal.RenderInternal) renderer.NodeRenderer {
	return &InlineRenderer{renderInternal: renderInternal}
}

func (r *InlineRenderer) renderInline(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(string(r.renderInternal.ProtectSafeAttrs(`<code class="language-math">`)))
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := util.EscapeHTML(segment.Value(source))
			if bytes.HasSuffix(value, []byte("\n")) {
				_, _ = w.Write(value[:len(value)-1])
				if c != n.LastChild() {
					_, _ = w.Write([]byte(" "))
				}
			} else {
				_, _ = w.Write(value)
			}
		}
		return ast.WalkSkipChildren, nil
	}
	_, _ = w.WriteString(`</code>`)
	return ast.WalkContinue, nil
}

// RegisterFuncs registers the renderer for inline math nodes
func (r *InlineRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindInline, r.renderInline)
}
