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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRenderConfig_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		expected *RenderConfig
		args     string
	}{
		{
			"empty", &RenderConfig{
				Meta: "table",
				Lang: "",
			}, "",
		},
		{
			"lang", &RenderConfig{
				Meta: "table",
				Lang: "test",
			}, "lang: test",
		},
		{
			"metatable", &RenderConfig{
				Meta: "table",
				Lang: "",
			}, "kmup: table",
		},
		{
			"metanone", &RenderConfig{
				Meta: "none",
				Lang: "",
			}, "kmup: none",
		},
		{
			"metadetails", &RenderConfig{
				Meta: "details",
				Lang: "",
			}, "kmup: details",
		},
		{
			"metawrong", &RenderConfig{
				Meta: "details",
				Lang: "",
			}, "kmup: wrong",
		},
		{
			"toc", &RenderConfig{
				TOC:  "true",
				Meta: "table",
				Lang: "",
			}, "include_toc: true",
		},
		{
			"tocfalse", &RenderConfig{
				TOC:  "false",
				Meta: "table",
				Lang: "",
			}, "include_toc: false",
		},
		{
			"toclang", &RenderConfig{
				Meta: "table",
				TOC:  "true",
				Lang: "testlang",
			}, `
				include_toc: true
				lang: testlang
				`,
		},
		{
			"complexlang", &RenderConfig{
				Meta: "table",
				Lang: "testlang",
			}, `
				kmup:
					lang: testlang
				`,
		},
		{
			"complexlang2", &RenderConfig{
				Meta: "table",
				Lang: "testlang",
			}, `
	lang: notright
	kmup:
		lang: testlang
`,
		},
		{
			"complexlang", &RenderConfig{
				Meta: "table",
				Lang: "testlang",
			}, `
	kmup:
		lang: testlang
`,
		},
		{
			"complex2", &RenderConfig{
				Lang: "two",
				Meta: "table",
				TOC:  "true",
			}, `
	lang: one
	include_toc: true
	kmup:
		details_icon: smiley
		meta: table
		include_toc: true
		lang: two
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &RenderConfig{
				Meta: "table",
				Lang: "",
			}
			err := yaml.Unmarshal([]byte(strings.ReplaceAll(tt.args, "\t", "    ")), got)
			require.NoError(t, err)

			assert.Equal(t, tt.expected.Meta, got.Meta)
			assert.Equal(t, tt.expected.Lang, got.Lang)
			assert.Equal(t, tt.expected.TOC, got.TOC)
		})
	}
}
