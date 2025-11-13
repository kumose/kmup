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
	"net/url"

	"github.com/kumose/kmup/modules/translation"

	"github.com/yuin/goldmark/ast"
)

// Header holds the data about a header.
type Header struct {
	Level int
	Text  string
	ID    string
}

func createTOCNode(toc []Header, lang string, detailsAttrs map[string]string) ast.Node {
	details := NewDetails()
	summary := NewSummary()

	for k, v := range detailsAttrs {
		details.SetAttributeString(k, []byte(v))
	}

	summary.AppendChild(summary, ast.NewString([]byte(translation.NewLocale(lang).TrString("toc"))))
	details.AppendChild(details, summary)
	ul := ast.NewList('-')
	details.AppendChild(details, ul)
	currentLevel := 6
	for _, header := range toc {
		if header.Level < currentLevel {
			currentLevel = header.Level
		}
	}
	for _, header := range toc {
		for currentLevel > header.Level {
			ul = ul.Parent().(*ast.List)
			currentLevel--
		}
		for currentLevel < header.Level {
			newL := ast.NewList('-')
			ul.AppendChild(ul, newL)
			currentLevel++
			ul = newL
		}
		li := ast.NewListItem(currentLevel * 2)
		a := ast.NewLink()
		a.Destination = []byte("#" + url.QueryEscape(header.ID))
		a.AppendChild(a, ast.NewString([]byte(header.Text)))
		li.AppendChild(li, a)
		ul.AppendChild(ul, li)
	}

	return details
}
