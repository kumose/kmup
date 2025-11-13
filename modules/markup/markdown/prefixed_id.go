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

package markdown

import (
	"bytes"
	"fmt"

	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/markup/common"
	"github.com/kumose/kmup/modules/util"

	"github.com/yuin/goldmark/ast"
)

type prefixedIDs struct {
	values container.Set[string]
}

// Generate generates a new element id.
func (p *prefixedIDs) Generate(value []byte, kind ast.NodeKind) []byte {
	dft := []byte("id")
	if kind == ast.KindHeading {
		dft = []byte("heading")
	}
	return p.GenerateWithDefault(value, dft)
}

// GenerateWithDefault generates a new element id.
func (p *prefixedIDs) GenerateWithDefault(value, dft []byte) []byte {
	result := common.CleanValue(value)
	if len(result) == 0 {
		result = dft
	}
	if !bytes.HasPrefix(result, []byte("user-content-")) {
		result = append([]byte("user-content-"), result...)
	}
	if p.values.Add(util.UnsafeBytesToString(result)) {
		return result
	}
	for i := 1; ; i++ {
		newResult := fmt.Sprintf("%s-%d", result, i)
		if p.values.Add(newResult) {
			return []byte(newResult)
		}
	}
}

// Put puts a given element id to the used ids table.
func (p *prefixedIDs) Put(value []byte) {
	p.values.Add(util.UnsafeBytesToString(value))
}

func newPrefixedIDs() *prefixedIDs {
	return &prefixedIDs{
		values: make(container.Set[string]),
	}
}
