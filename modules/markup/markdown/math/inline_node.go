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
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// Inline struct represents inline math e.g. $...$ or \(...\)
type Inline struct {
	ast.BaseInline
}

// Inline implements Inline.Inline.
func (n *Inline) Inline() {}

// IsBlank returns if this inline node is empty
func (n *Inline) IsBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*ast.Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

// Dump renders this inline math as debug
func (n *Inline) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindInline is the kind for math inline
var KindInline = ast.NewNodeKind("MathInline")

// Kind returns KindInline
func (n *Inline) Kind() ast.NodeKind {
	return KindInline
}

// NewInline creates a new ast math inline node
func NewInline() *Inline {
	return &Inline{
		BaseInline: ast.BaseInline{},
	}
}
