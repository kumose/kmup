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

package console

import (
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/typesniffer"
	"github.com/kumose/kmup/modules/util"

	trend "github.com/buildkite/terminal-to-html/v3"
)

func init() {
	markup.RegisterRenderer(Renderer{})
}

// Renderer implements markup.Renderer
type Renderer struct{}

var _ markup.RendererContentDetector = (*Renderer)(nil)

// Name implements markup.Renderer
func (Renderer) Name() string {
	return "console"
}

// Extensions implements markup.Renderer
func (Renderer) Extensions() []string {
	return []string{".sh-session"}
}

// SanitizerRules implements markup.Renderer
func (Renderer) SanitizerRules() []setting.MarkupSanitizerRule {
	return []setting.MarkupSanitizerRule{
		{Element: "span", AllowAttr: "class", Regexp: `^term-((fg[ix]?|bg)\d+|container)$`},
	}
}

// CanRender implements markup.RendererContentDetector
func (Renderer) CanRender(filename string, sniffedType typesniffer.SniffedType, prefetchBuf []byte) bool {
	if !sniffedType.IsTextPlain() {
		return false
	}

	s := util.UnsafeBytesToString(prefetchBuf)
	rs := []rune(s)
	cnt := 0
	firstErrPos := -1
	isCtrlSep := func(p int) bool {
		return p < len(rs) && (rs[p] == ';' || rs[p] == 'm')
	}
	for i, c := range rs {
		if c == 0 {
			return false
		}
		if c == '\x1b' {
			match := i+1 < len(rs) && rs[i+1] == '['
			if match && (isCtrlSep(i+2) || isCtrlSep(i+3) || isCtrlSep(i+4) || isCtrlSep(i+5)) {
				cnt++
			}
		}
		if c == utf8.RuneError && firstErrPos == -1 {
			firstErrPos = i
		}
	}
	if firstErrPos != -1 && firstErrPos != len(rs)-1 {
		return false
	}
	return cnt >= 2 // only render it as console output if there are at least two escape sequences
}

// Render renders terminal colors to HTML with all specific handling stuff.
func (Renderer) Render(ctx *markup.RenderContext, input io.Reader, output io.Writer) error {
	buf, err := io.ReadAll(input)
	if err != nil {
		return err
	}
	buf = []byte(trend.Render(buf))
	buf = bytes.ReplaceAll(buf, []byte("\n"), []byte(`<br>`))
	_, err = output.Write(buf)
	return err
}
