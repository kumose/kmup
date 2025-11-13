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
	"html/template"

	"github.com/kumose/kmup/modules/markup/internal"
	kmupUtil "github.com/kumose/kmup/modules/util"

	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Block render output:
// 	<pre class="code-block is-loading"><code class="language-math display">...</code></pre>
//
// Keep in mind that there is another "code block" render in "func (r *GlodmarkRender) highlightingRenderer"
// "highlightingRenderer" outputs the math block with extra "chroma" class:
// 	<pre class="code-block is-loading"><code class="chroma language-math display">...</code></pre>
//
// Special classes:
// * "is-loading": show a loading indicator
// * "display": used by JS to decide to render as a block, otherwise render as inline

// BlockRenderer represents a renderer for math Blocks
type BlockRenderer struct {
	renderInternal *internal.RenderInternal
}

// NewBlockRenderer creates a new renderer for math Blocks
func NewBlockRenderer(renderInternal *internal.RenderInternal) renderer.NodeRenderer {
	return &BlockRenderer{renderInternal: renderInternal}
}

// RegisterFuncs registers the renderer for math Blocks
func (r *BlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindBlock, r.renderBlock)
}

func (r *BlockRenderer) writeLines(w util.BufWriter, source []byte, n gast.Node) {
	l := n.Lines().Len()
	for i := range l {
		line := n.Lines().At(i)
		_, _ = w.Write(util.EscapeHTML(line.Value(source)))
	}
}

func (r *BlockRenderer) renderBlock(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*Block)
	if entering {
		codeHTML := kmupUtil.Iif[template.HTML](n.Inline, "", `<pre class="code-block is-loading">`) + `<code class="language-math display">`
		_, _ = w.WriteString(string(r.renderInternal.ProtectSafeAttrs(codeHTML)))
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString(`</code>` + kmupUtil.Iif(n.Inline, "", `</pre>`) + "\n")
	}
	return gast.WalkContinue, nil
}
