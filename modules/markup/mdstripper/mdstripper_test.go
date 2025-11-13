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

package mdstripper

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownStripper(t *testing.T) {
	type testItem struct {
		markdown      string
		expectedText  []string
		expectedLinks []string
	}

	list := []testItem{
		{
			`
## This is a title

This is [one](link) to paradise.
This **is emphasized**.
This: should coalesce.

` + "```" + `
This is a code block.
This should not appear in the output at all.
` + "```" + `

* Bullet 1
* Bullet 2

A HIDDEN ` + "`" + `GHOST` + "`" + ` IN THIS LINE.
		`,
			[]string{
				"This is a title",
				"This is",
				"to paradise.",
				"This",
				"is emphasized",
				".",
				"This: should coalesce.",
				"Bullet 1",
				"Bullet 2",
				"A HIDDEN",
				"IN THIS LINE.",
			},
			[]string{
				"link",
			},
		},
		{
			"Simply closes: #29 yes",
			[]string{
				"Simply closes: #29 yes",
			},
			[]string{},
		},
		{
			"Simply closes: !29 yes",
			[]string{
				"Simply closes: !29 yes",
			},
			[]string{},
		},
	}

	for _, test := range list {
		text, links := StripMarkdown([]byte(test.markdown))
		rawlines := strings.Split(text, "\n")
		lines := make([]string, 0, len(rawlines))
		for _, line := range rawlines {
			line := strings.TrimSpace(line)
			if line != "" {
				lines = append(lines, line)
			}
		}
		assert.Equal(t, test.expectedText, lines)
		assert.Equal(t, test.expectedLinks, links)
	}
}
