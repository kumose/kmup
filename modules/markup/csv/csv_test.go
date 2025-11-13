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

package markup

import (
	"strings"
	"testing"

	"github.com/kumose/kmup/modules/markup"

	"github.com/stretchr/testify/assert"
)

func TestRenderCSV(t *testing.T) {
	var render Renderer
	kases := map[string]string{
		"a":        "<table class=\"data-table\"><tr><th class=\"line-num\">1</th><th>a</th></tr></table>",
		"1,2":      "<table class=\"data-table\"><tr><th class=\"line-num\">1</th><th>1</th><th>2</th></tr></table>",
		"1;2\n3;4": "<table class=\"data-table\"><tr><th class=\"line-num\">1</th><th>1</th><th>2</th></tr><tr><td class=\"line-num\">2</td><td>3</td><td>4</td></tr></table>",
		"<br/>":    "<table class=\"data-table\"><tr><th class=\"line-num\">1</th><th>&lt;br/&gt;</th></tr></table>",
	}

	for k, v := range kases {
		var buf strings.Builder
		err := render.Render(markup.NewRenderContext(t.Context()), strings.NewReader(k), &buf)
		assert.NoError(t, err)
		assert.Equal(t, v, buf.String())
	}
}
