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

import "github.com/yuin/goldmark/ast"

// Block represents a display math block e.g. $$...$$ or \[...\]
type Block struct {
	ast.BaseBlock
	Dollars bool
	Indent  int
	Closed  bool
	Inline  bool
}

// KindBlock is the node kind for math blocks
var KindBlock = ast.NewNodeKind("MathBlock")

// NewBlock creates a new math Block
func NewBlock(dollars bool, indent int) *Block {
	return &Block{
		Dollars: dollars,
		Indent:  indent,
	}
}

// Dump dumps the block to a string
func (n *Block) Dump(source []byte, level int) {
	m := map[string]string{}
	ast.DumpHelper(n, source, level, m, nil)
}

// Kind returns KindBlock for math Blocks
func (n *Block) Kind() ast.NodeKind {
	return KindBlock
}

// IsRaw returns true as this block should not be processed further
func (n *Block) IsRaw() bool {
	return true
}
