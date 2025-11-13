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

package fileicon

import (
	"html/template"
	"strings"

	"github.com/kumose/kmup/modules/setting"
)

type RenderedIconPool struct {
	IconSVGs map[string]template.HTML
}

func NewRenderedIconPool() *RenderedIconPool {
	return &RenderedIconPool{
		IconSVGs: make(map[string]template.HTML),
	}
}

func (p *RenderedIconPool) RenderToHTML() template.HTML {
	if len(p.IconSVGs) == 0 {
		return ""
	}
	sb := &strings.Builder{}
	sb.WriteString(`<div class=tw-hidden>`)
	for _, icon := range p.IconSVGs {
		sb.WriteString(string(icon))
	}
	sb.WriteString(`</div>`)
	return template.HTML(sb.String())
}

func RenderEntryIconHTML(renderedIconPool *RenderedIconPool, entry *EntryInfo) template.HTML {
	if setting.UI.FileIconTheme == "material" {
		return DefaultMaterialIconProvider().EntryIconHTML(renderedIconPool, entry)
	}
	return BasicEntryIconHTML(entry)
}
