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
	"github.com/kumose/kmup/modules/markup/internal"
	kmupUtil "github.com/kumose/kmup/modules/util"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Options struct {
	Enabled                  bool
	ParseInlineDollar        bool // inline $$ xxx $$ text
	ParseInlineParentheses   bool // inline \( xxx \) text
	ParseBlockDollar         bool // block $$ multiple-line $$ text
	ParseBlockSquareBrackets bool // block \[ multiple-line \] text
}

// Extension is a math extension
type Extension struct {
	renderInternal *internal.RenderInternal
	options        Options
}

// NewExtension creates a new math extension with the provided options
func NewExtension(renderInternal *internal.RenderInternal, opts ...Options) *Extension {
	opt := kmupUtil.OptionalArg(opts)
	r := &Extension{
		renderInternal: renderInternal,
		options:        opt,
	}
	return r
}

// Extend extends goldmark with our parsers and renderers
func (e *Extension) Extend(m goldmark.Markdown) {
	if !e.options.Enabled {
		return
	}

	var inlines []util.PrioritizedValue
	if e.options.ParseInlineParentheses {
		inlines = append(inlines, util.Prioritized(NewInlineParenthesesParser(), 501))
	}
	inlines = append(inlines, util.Prioritized(NewInlineDollarParser(e.options.ParseInlineDollar), 502))

	m.Parser().AddOptions(parser.WithInlineParsers(inlines...))
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewBlockParser(e.options.ParseBlockDollar, e.options.ParseBlockSquareBrackets), 701),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewBlockRenderer(e.renderInternal), 501),
		util.Prioritized(NewInlineRenderer(e.renderInternal), 502),
	))
}
