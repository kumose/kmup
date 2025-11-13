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
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/typesniffer"

	"github.com/stretchr/testify/assert"
)

func TestRenderConsole(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"\x1b[37m\x1b[40mnpm\x1b[0m \x1b[0m\x1b[32minfo\x1b[0m \x1b[0m\x1b[35mit worked if it ends with\x1b[0m ok", `<span class="term-fg37 term-bg40">npm</span> <span class="term-fg32">info</span> <span class="term-fg35">it worked if it ends with</span> ok`},
		{"\x1b[1;2m \x1b[123m 啊", `<span class="term-fg2">  啊</span>`},
		{"\x1b[1;2m \x1b[123m \xef", `<span class="term-fg2">  �</span>`},
		{"\x1b[1;2m \x1b[123m \xef \xef", ``},
		{"\x1b[12", ``},
		{"\x1b[1", ``},
		{"\x1b[FOO\x1b[", ``},
		{"\x1b[mFOO\x1b[m", `FOO`},
	}

	var render Renderer
	for i, c := range cases {
		var buf strings.Builder
		st := typesniffer.DetectContentType([]byte(c.input))
		canRender := render.CanRender("test", st, []byte(c.input))
		if c.expected == "" {
			assert.False(t, canRender, "case %d: expected not to render", i)
			continue
		}

		assert.True(t, canRender)
		err := render.Render(markup.NewRenderContext(t.Context()), strings.NewReader(c.input), &buf)
		assert.NoError(t, err)
		assert.Equal(t, c.expected, buf.String())
	}
}
