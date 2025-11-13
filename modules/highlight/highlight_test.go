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

package highlight

import (
	"html/template"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func lines(s string) (out []template.HTML) {
	// "" => [], "a" => ["a"], "a\n" => ["a\n"], "a\nb" => ["a\n", "b"] (each line always includes EOL "\n" if it exists)
	out = make([]template.HTML, 0)
	s = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(s), "\n", ""), `\n`, "\n")
	for {
		if p := strings.IndexByte(s, '\n'); p != -1 {
			out = append(out, template.HTML(s[:p+1]))
			s = s[p+1:]
		} else {
			break
		}
	}
	if s != "" {
		out = append(out, template.HTML(s))
	}
	return out
}

func TestFile(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		want      []template.HTML
		lexerName string
	}{
		{
			name:      "empty.py",
			code:      "",
			want:      lines(""),
			lexerName: "Python",
		},
		{
			name:      "empty.js",
			code:      "",
			want:      lines(""),
			lexerName: "JavaScript",
		},
		{
			name:      "empty.yaml",
			code:      "",
			want:      lines(""),
			lexerName: "YAML",
		},
		{
			name:      "tags.txt",
			code:      "<>",
			want:      lines("&lt;&gt;"),
			lexerName: "Plaintext",
		},
		{
			name:      "tags.py",
			code:      "<>",
			want:      lines(`<span class="o">&lt;</span><span class="o">&gt;</span>`),
			lexerName: "Python",
		},
		{
			name:      "eol-no.py",
			code:      "a=1",
			want:      lines(`<span class="n">a</span><span class="o">=</span><span class="mi">1</span>`),
			lexerName: "Python",
		},
		{
			name:      "eol-newline1.py",
			code:      "a=1\n",
			want:      lines(`<span class="n">a</span><span class="o">=</span><span class="mi">1</span>\n`),
			lexerName: "Python",
		},
		{
			name: "eol-newline2.py",
			code: "a=1\n\n",
			want: lines(`
<span class="n">a</span><span class="o">=</span><span class="mi">1</span>\n
\n
			`,
			),
			lexerName: "Python",
		},
		{
			name: "empty-line-with-space.py",
			code: strings.ReplaceAll(strings.TrimSpace(`
def:
    a=1

b=''
{space}
c=2
			`), "{space}", "    "),
			want: lines(`
<span class="n">def</span><span class="p">:</span>\n
    <span class="n">a</span><span class="o">=</span><span class="mi">1</span>\n
\n
<span class="n">b</span><span class="o">=</span><span class="sa"></span><span class="s1">&#39;</span><span class="s1">&#39;</span>\n
    \n
<span class="n">c</span><span class="o">=</span><span class="mi">2</span>`,
			),
			lexerName: "Python",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, lexerName, err := File(tt.name, "", []byte(tt.code))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, out)
			assert.Equal(t, tt.lexerName, lexerName)
		})
	}
}

func TestPlainText(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []template.HTML
	}{
		{
			name: "empty.py",
			code: "",
			want: lines(""),
		},
		{
			name: "tags.py",
			code: "<>",
			want: lines("&lt;&gt;"),
		},
		{
			name: "eol-no.py",
			code: "a=1",
			want: lines(`a=1`),
		},
		{
			name: "eol-newline1.py",
			code: "a=1\n",
			want: lines(`a=1\n`),
		},
		{
			name: "eol-newline2.py",
			code: "a=1\n\n",
			want: lines(`
a=1\n
\n
			`),
		},
		{
			name: "empty-line-with-space.py",
			code: strings.ReplaceAll(strings.TrimSpace(`
def:
    a=1

b=''
{space}
c=2
			`), "{space}", "    "),
			want: lines(`
def:\n
    a=1\n
\n
b=&#39;&#39;\n
    \n
c=2`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := PlainText([]byte(tt.code))
			assert.Equal(t, tt.want, out)
		})
	}
}
