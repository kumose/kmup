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

package asciicast

import (
	"fmt"
	"io"
	"net/url"

	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/setting"
)

func init() {
	markup.RegisterRenderer(Renderer{})
}

// Renderer implements markup.Renderer for asciicast files.
// See https://github.com/asciinema/asciinema/blob/develop/doc/asciicast-v2.md
type Renderer struct{}

// Name implements markup.Renderer
func (Renderer) Name() string {
	return "asciicast"
}

// Extensions implements markup.Renderer
func (Renderer) Extensions() []string {
	return []string{".cast"}
}

const (
	playerClassName = "asciinema-player-container"
	playerSrcAttr   = "data-asciinema-player-src"
)

// SanitizerRules implements markup.Renderer
func (Renderer) SanitizerRules() []setting.MarkupSanitizerRule {
	return []setting.MarkupSanitizerRule{{Element: "div", AllowAttr: playerSrcAttr}}
}

// Render implements markup.Renderer
func (Renderer) Render(ctx *markup.RenderContext, _ io.Reader, output io.Writer) error {
	rawURL := fmt.Sprintf("%s/%s/%s/raw/%s/%s",
		setting.AppSubURL,
		url.PathEscape(ctx.RenderOptions.Metas["user"]),
		url.PathEscape(ctx.RenderOptions.Metas["repo"]),
		ctx.RenderOptions.Metas["RefTypeNameSubURL"],
		url.PathEscape(ctx.RenderOptions.RelativePath),
	)
	return ctx.RenderInternal.FormatWithSafeAttrs(output, `<div class="%s" %s="%s"></div>`, playerClassName, playerSrcAttr, rawURL)
}
