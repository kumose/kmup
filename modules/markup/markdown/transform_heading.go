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
	"fmt"

	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/util"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func (g *ASTTransformer) transformHeading(_ *markup.RenderContext, v *ast.Heading, reader text.Reader, tocList *[]Header) {
	for _, attr := range v.Attributes() {
		if _, ok := attr.Value.([]byte); !ok {
			v.SetAttribute(attr.Name, fmt.Appendf(nil, "%v", attr.Value))
		}
	}
	txt := v.Text(reader.Source()) //nolint:staticcheck // Text is deprecated
	header := Header{
		Text:  util.UnsafeBytesToString(txt),
		Level: v.Level,
	}
	if id, found := v.AttributeString("id"); found {
		header.ID = util.UnsafeBytesToString(id.([]byte))
	}
	*tocList = append(*tocList, header)
	g.applyElementDir(v)
}
